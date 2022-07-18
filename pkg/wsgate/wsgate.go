package wsgate

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	spec "github.com/IBM/workload-security-guard/pkg/apis/wsecurity/v1"

	"github.com/IBM/go-security-plugs/iodup"
	"github.com/IBM/go-security-plugs/iofilter"

	pi "github.com/IBM/go-security-plugs/pluginterfaces"

	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

const version string = "0.0.1"
const name string = "wsgate"

//type StateReq time.Time
//type StateResp time.Time

type ctxKey string

type plug struct {
	name    string
	version string

	kubemgr Kubemgr

	// Add here any other state the extension needs
	guardUrl  string
	namespace string
	serviceId string
	cmname    bool
	wsGate    *spec.GuardianSpec
	//blocked             []string
	//numOk               uint32
	//lastConsultReported time.Time
	//numConsultsCount    uint16
	httpc http.Client
	cycle int
	//allowedPile         spec.Pile
	pile       spec.Pile
	index      uint32
	ongoing    map[uint32]*spec.SessionProfile
	statistics map[string]uint32
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))[0:16]
}

func (p *plug) Shutdown() {
	pi.Log.Infof("%s: Shutdown", p.name)
}

func (p *plug) PlugName() string {
	return p.name
}

func (p *plug) PlugVersion() string {
	return p.version
}

func ReadUserIP(req *http.Request) string {
	IPAddress := req.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = req.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = req.RemoteAddr
	}
	return IPAddress
}

func (p *plug) screenRequestBody(req *http.Request, bp *spec.BodyProfile) error {
	var decission string
	testBodyHist := true
	maxBody := int64(1000000)
	doneBody := false

	if req.ContentLength < maxBody {
		// TBD - validate content-type params returned by ParseMediaType!
		ctype, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
		if err != nil {
			ctype = "application/octet-stream"
		}
		method := req.Method
		pi.Log.Debugf("%s Request Content-Type: %s  (params %v)\n", p.name, ctype, params)
		pi.Log.Debugf("%s Request Method: %s\n", p.name, method)
		if ctype == "application/json" {
			pi.Log.Debugf("Processing json!\n")
			doneBody = true
			dup := iodup.New(req.Body, 2, 128, 8192)

			var j interface{}
			dec := json.NewDecoder(&dup.Output[1])
			err = dec.Decode(&j)
			if err != nil {
				pi.Log.Debugf("Error Decoding json! %v\n", err)
			}
			pi.Log.Debugf("Json! %v\n", j)
			bp.Structured = new(spec.StructuredProfile)
			bp.Structured.Profile(j)
			pi.Log.Debugf("JsonProfile: %s\n", bp.Structured.Marshal(0))
			req.Body = &dup.Output[0]

		} else if method == "POST" || method == "PUT" || method == "PATCH" {
			pi.Log.Debugf("Processing POST/PUT/PATCH\n")
			if ctype == "application/x-www-form-urlencoded" {
				pi.Log.Debugf("Processing application/x-www-form-urlencoded!\n")
				doneBody = true
				// Option A
				//body, err := ioutil.ReadAll(req.Body)
				//if err != nil {
				//	return nil, errors.New("error reading body")
				//}
				// Option B
				dup := iodup.New(req.Body, 2, 128, 8192)

				// option 1
				//req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
				req.Body = &dup.Output[1]
				req.ParseForm()
				pi.Log.Debugf("Form Data %s\n", req.Form)

				bp.Structured = new(spec.StructuredProfile)
				bp.Structured.ProfilePostForm(req.PostForm)
				s := bp.Structured.Marshal(0)
				pi.Log.Debugf("--> %d %s\n", len(s), s)
				pi.Log.Debugf("JsonProfile: %s\n", s)
				// option 2
				// create new request for parsing the body
				//reqCopy, _ := http.NewRequest(req.Method, req.URL.String(), bytes.NewReader(body))
				//reqCopy, _ := http.NewRequest(req.Method, req.URL.String(), &dup.Output[1])
				//reqCopy.Header = req.Header
				//reqCopy.ParseForm()
				//pi.Log.Debugf("Form Data %s\n", reqCopy.Form)

				// reset orig request
				//req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
				req.Body = &dup.Output[0]
			} else if ctype == "multipart/form-data" {
				pi.Log.Debugf("Processing multipart/form-data!\n")
				doneBody = true

				// Option A
				//body, err := ioutil.ReadAll(req.Body)
				//if err != nil {
				//	return nil, errors.New("error reading body")
				//}
				// Option B
				dup := iodup.New(req.Body, 2, 128, 8192)

				// option 1
				//req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
				req.Body = &dup.Output[1]
				req.ParseMultipartForm(maxBody)
				pi.Log.Debugf("Form Data %s\n", req.Form)

				bp.Structured = new(spec.StructuredProfile)
				bp.Structured.ProfilePostForm(req.PostForm)
				//s := bp.Structured.Marshal(0)
				//pi.Log.Debugf("--> %d %s\n", len(s), s)
				//pi.Log.Debugf("JsonProfile: %s\n", s)
				// option 2
				// create new request for parsing the body
				//reqCopy, _ := http.NewRequest(req.Method, req.URL.String(), bytes.NewReader(body))
				//reqCopy, _ := http.NewRequest(req.Method, req.URL.String(), &dup.Output[1])
				//reqCopy.Header = req.Header
				//reqCopy.ParseForm()
				//pi.Log.Debugf("Form Data %s\n", reqCopy.Form)

				// reset orig request
				//req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
				req.Body = &dup.Output[0]
			}
		}
		if !doneBody && testBodyHist {
			fmt.Printf("Analyze Body\n")

			// Asynchrniously stream bytes from the original resp.Body
			// to a new resp.Body
			//req.Body = iofilter.New(req.Body, requestFilter, sp.ReqBody)
		}
	}
	ctrl := p.wsGate.Control
	//fmt.Printf("screenBodyRequest ctrl %v\n", ctrl)
	var criteria *spec.Critiria
	if criteria = p.wsGate.Configured; ctrl.Auto {
		criteria = p.wsGate.Learned
	}
	if (criteria != nil) && criteria.Active {
		decission = criteria.ReqBody.Decide(bp)
	}
	if ctrl.Alert {
		if decission != "" {
			pi.Log.Infof("Alert HttpRequestBody: %s", decission)
			p.stats("ReqBodyNok")
			return errors.New(decission)
		}
	}

	p.stats("ReqBodyOk")
	return nil
}

func (p *plug) screenRequest(req *http.Request, rp *spec.ReqProfile) error {
	var decission string

	pi.Log.Debugf("screenRequest: ctrl %v", p.wsGate.Control)
	p.stats("Total")
	// Request client and server identities
	cip, cport, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		decission += fmt.Sprintf("illegal req.RemoteAddr %s", err.Error())
	}
	sip, sport, err := net.SplitHostPort(req.URL.Host)

	if err != nil {
		decission += fmt.Sprintf("illegal req.URL.Host %s", err.Error())
	}
	pi.Log.Debugf("Client: %s port %s", cip, cport)
	pi.Log.Debugf("Server: %s port %s", sip, sport)

	// Request principles
	pi.Log.Debugf("req.Method %s", req.Method)
	pi.Log.Debugf("req.Proto %s", req.Proto)
	pi.Log.Debugf("scheme: %s", req.URL.Scheme)
	pi.Log.Debugf("opaque: %s", req.URL.Opaque)

	pi.Log.Debugf("ContentLength: %d", req.ContentLength)
	pi.Log.Debugf("Trailer: %#v", req.Trailer)

	// TBD req.Form

	ip := net.ParseIP(cip)
	rp.Profile(req, ip)
	//fmt.Println(rp.Marshal(0))
	//fmt.Println(p.wsGate.Configured)
	ctrl := p.wsGate.Control
	//fmt.Printf("screenRequest ctrl %v\n", ctrl)
	var criteria *spec.Critiria
	fmt.Printf("p.wsGate %v ctrl %v\n", p.wsGate, ctrl)
	if criteria = p.wsGate.Configured; ctrl.Auto {
		pi.Log.Infof("Using Learned Critiria!")
		criteria = p.wsGate.Learned
	}
	pi.Log.Infof("Critiria! %v", criteria)
	if criteria != nil {
		pi.Log.Infof("Critiria.Active! %v", criteria.Active)
	}
	if (criteria != nil) && criteria.Active {
		pi.Log.Infof("Deciding using Critiria!")
		decission = criteria.Req.Decide(rp)
	}
	if ctrl.Alert {
		if decission != "" {
			pi.Log.Infof("Alert HttpRequest: %s", decission)
			p.stats("ReqNok")
			return errors.New(decission)
		}
	}

	p.stats("ReqOk")
	return nil
	/*
		if decission != "" {
			// potentially consult guard before rejecting
			pi.Log.Infof("Guardian refused to allow: %s", decission)
			if ctrl.Consult {
				minuete := time.Now().Truncate(time.Minute)
				if p.lastConsultReported != minuete {
					p.lastConsultReported = minuete
					p.numConsultsCount = 0
				}
				if p.numConsultsCount < ctrl.RequestsPerMinuete {
					p.numConsultsCount = p.numConsultsCount + 1
					pi.Log.Infof("Consulting Guard %d/%d", p.numConsultsCount, ctrl.RequestsPerMinuete)
					decission = p.consultOnRequest(rp)
					//pi.Log.Infof("Guard said: %s", decission)
				}
			}
		}
	*/

	/*
		//decoded path
		path := req.URL.Path
		pathProfile := p.gateConfig.ProfilePath(path)
		pi.Log.Infof("path profile %v", pathProfile)

		//decoded query string
		query := req.URL.Query()
		queryProfile := p.gateConfig.ProfileQueryString(query)
		pi.Log.Infof("query: %#v", queryProfile)

		//http headers
		hProfile := p.gateConfig.ProfileHttpHeaders(req.Header)

		pi.Log.Infof("Headers: %#v", hProfile)

		//http Trailers
		trailerProfile := spec.ProfileKeyVals(req.Trailer)
		pi.Log.Infof("Trailer: %#v", trailerProfile)
	*/
	/*
		// fingerprints representing the sender of the request and the event to be processed
		symbols := make([]string, 0, 12)
		symbols = append(symbols,
			req.Method,
			req.Proto,
			u.Scheme,
			opaque,

		// fingerprints representing the sender of the request and the event to be processed
		fingerprints := make([]string, 0, 12)
		fingerprints = append(fingerprints,
			pathSplits[0],
			GetMD5Hash(queryKeys),
			headers.Get("Transfer-Encoding"),
			headers.Get("Content-Encoding"),
			headers.Get("Keep-Alive"),
			headers.Get("Connection"),
			headers.Get("X-Forwarded-For"),
			headers.Get("Cache-Control"),
			headers.Get("Via"),
			acceptHeaderVals,
			contentHeaderVals,
			userAgentVals,
			allHeaderKeys,
			httpinfo["protocol"])
		pi.Log(fingerprints)
		for i, val := range fingerprints {
			fingerprints[i] = GetMD5Hash(val)
		}
	*/
	/*


		contentEncoding := r.Header.Get("content-encoding")
		transferEncoding := r.Header.Get("transfer-encoding")
		keepAlive := r.Header.Get("keep-alive")
		connection := r.Header.Get("Connection")
		xForwardedFor := r.Header.Get("x-forwarded-for")
		cacheControl := r.Header.Get("cache-control")
		via := r.Header.Get("via")

		log.Info("DH> userAgentVals ", userAgentVals)
		log.Info("DH> contentEncoding ", contentEncoding)
		log.Info("DH> transferEncoding ", transferEncoding)
		log.Info("DH> keepAlive ", keepAlive)
		log.Info("DH> connection ", connection)
		log.Info("DH> xForwardedFor ", xForwardedFor)
		log.Info("DH> cacheControl ", cacheControl)
		log.Info("DH> via ", via)
	*/
	//var d = new Date();
	//h := make(map[string]string)

	//markers := make([]float32, 0, 12)
	//integers := make([]int, 0, 12)
	//roundedMarkers := make([]float32, 0, 12)
	//histograms := make([][]int, 0, 12)

	// Create a sorted slice of all header keys

	// Create a sorted slice of all query leys

	/*
		roundedMarkers.append(fingerprints, d.getDay()/6)
		roundedMarkers.append(fingerprints, d.getHours()/23)
		console.log(roundedMarkers)

		console.log(httpreq.body)
		console.log(otherHeaderVals)
		console.log(queryContent)


		integers.append(integers, parseInt(httpreq.size)) // Content-Length  - size of body
		integers.append(integers, otherHeaderVals.length)
		integers.append(integers, queryContent.length)
		integers.append(integers, cookieVals.length)
		integers.append(integers, pathSplits[0].length)
		integers.append(integers, allHeaderVals.length)
		console.log(markers, integers)



		histograms.append(histograms, hist(httpreq.body))
		histograms.append(histograms, hist(otherHeaderVals))
		histograms.append(histograms, hist(queryContent))
		histograms.append(histograms, hist(cookieVals))
		histograms.append(histograms, hist(allHeaderVals))
		console.log(histograms)

		fingerprint_path= pathSplits[0]


		var triggerInstance = headers["x-request-id"]||uuid.v4()




		const dataout = JSON.stringify({
					gateId:   gate
				, serviceId: unit
				, triggerInstance: triggerInstance
				, data: {
						fingerprints: fingerprints
					, markers: markers
					, integers: integers
					, roundedMarkers: roundedMarkers
					, histograms: histograms
				}
			});

		console.log(unit, dataout);
		postRequest("Path: "+fingerprint_path, "/eval", dataout, callback)
	*/
}

func (p *plug) screenResponse(resp *http.Response, rp *spec.RespProfile) error {
	rp.Profile(resp)
	//fmt.Println(rp.Marshal(0))
	//fmt.Println(p.wsGate.Configured)
	ctrl := p.wsGate.Control
	var decission string
	var criteria *spec.Critiria
	if criteria = p.wsGate.Configured; ctrl.Auto {
		criteria = p.wsGate.Learned
	}
	if (criteria != nil) && criteria.Active {
		decission = criteria.Resp.Decide(rp)
	}
	if ctrl.Alert {
		if decission != "" {
			pi.Log.Infof("Alert HttpResponse: %s", decission)
			p.stats("RespNok")
			return errors.New(decission)
		}
	}
	p.stats("RespOk")
	return nil
}

func responseFilter(buf []byte, state interface{}) {
	bp := state.(spec.BodyProfile)
	if bp.Unstructured == nil {
		bp.Unstructured = new(spec.SimpleValProfile)
	}
	h := make([]int, 8)

	for _, c := range buf {
		switch {
		case (c >= 97 && c <= 122) || (c >= 48 && c <= 57) || (c == 32): //a..z, 0..9, <SPACE>
			h[0]++
		case c >= 127 || c <= 31: // All non ascii unicodes, ascii 0-31, <DEL>
			h[1]++
		case c == 34 || c == 96 || c == 39: // ascii quatations  - TBD IF NEED TO BE extended with other suspects
			h[2]++
		case c == 59: // ; - TBD IF NEED TO BE extended with other suspects
			h[3]++
		default: // anything else: !#$%&()*+,-./:<=>?@[\]^_{|}~
			h[7]++
		}
	}
	//fmt.Printf("responseFilter Histogram: %v\n", h)
}

/*
func requestJsonFilter(buf []byte, state interface{}) {
	var m map[string]interface{} // declaring a map for key names as string and values as interface
	_ = json.Unmarshal(buf, &m)

	//fmt.Printf("requestJsonFilter: %v\n", m)
}

func requestFilter(buf []byte, state interface{}) {
	h := make([]int, 8)

	for _, c := range buf {
		switch {
		case (c >= 97 && c <= 122) || (c >= 48 && c <= 57) || (c == 32): //a..z, 0..9, <SPACE>
			h[0]++
		case c >= 127 || c <= 31: // All non ascii unicodes, ascii 0-31, <DEL>
			h[1]++
		case c == 34 || c == 96 || c == 39: // ascii quatations  - TBD IF NEED TO BE extended with other suspects
			h[2]++
		case c == 59: // ; - TBD IF NEED TO BE extended with other suspects
			h[3]++
		default: // anything else: !#$%&()*+,-./:<=>?@[\]^_{|}~
			h[7]++

		}
	}
	//fmt.Printf("requestFilter Histogram: %v\n", h)
}
*/
func (p *plug) periodical(ctx context.Context) bool {
	// Find the session
	index, okIndex := ctx.Value(ctxKey("ReqIndex")).(uint32)
	if !okIndex { // This should never happen!
		pi.Log.Warnf("%s ........... Periodical missing context!", p.name)
		return false
	}
	sp, spExists := p.ongoing[index]
	if !spExists { // This should never happen!
		pi.Log.Warnf("%s ........... Periodical missing session!", p.name)
		return false
	}

	now := time.Now()
	pp := &sp.Process
	//fmt.Printf("sp.ReqTime %v now %v\n", sp.ReqTime, now)
	pp.Profile(sp.ReqTime, sp.ReqTime, now)
	//fmt.Println(pp.Marshal(0))

	ctrl := p.wsGate.Control
	var decission string
	var criteria *spec.Critiria
	if criteria = p.wsGate.Configured; ctrl.Auto {
		criteria = p.wsGate.Learned
	}
	if (criteria != nil) && criteria.Active {
		decission = criteria.Process.Decide(pp)
	}

	if ctrl.Alert {
		if decission != "" {
			pi.Log.Infof("Alert during periodical: %s", decission)
			p.stats("PeriodicalNok")
			return true
		}
	}
	p.stats("PeriodicalOk")
	return false
}

func (p *plug) ApproveRequest(req *http.Request) (*http.Request, error) {

	pi.Log.Debugf("%s: ApproveRequest started", p.name)

	//if req.Header.Get("X-Block-Req") != "" {
	//	pi.Log.Infof("%s ........... Blocked During Request! returning an error!", p.name)
	//	return nil, errors.New("request blocked")
	//}

	//for name, values := range req.Header {
	// Loop over all values for the name.
	//	for _, value := range values {
	//pi.Log.Debugf("%s Request Header: %s: %s", p.name, name, value)
	//	}
	//}
	pi.Log.Debugf("%s Request Content-Length: %d\n", p.name, req.ContentLength)

	// TBD for dev stage only - read confi with every request
	// replace with a mechansim that reads config only if X min passed and the pod is still up
	p.wsGate = p.kubemgr.FetchConfig(p.namespace, p.serviceId, p.cmname)
	sp := new(spec.SessionProfile)
	sp.ReqTime = time.Now()

	// Req
	if p.screenRequest(req, &sp.Req) != nil {
		if p.wsGate.Control.Block {
			return nil, errors.New("secuirty blocked")
		} else {
			sp.Alert = true
		}
	}

	// Req Body
	if p.screenRequestBody(req, &sp.ReqBody) != nil {
		if p.wsGate.Control.Block {
			return nil, errors.New("secuirty blocked")
		} else {
			sp.Alert = true
		}
	}

	p.ongoing[p.index] = sp

	newCtx, cancelFunction := context.WithCancel(req.Context())
	newCtx = context.WithValue(newCtx, ctxKey("ReqIndex"), p.index)
	req = req.WithContext(newCtx)

	if p.periodical(newCtx) {
		pi.Log.Infof("Blocked on periodical during reuqest!")
		if p.wsGate.Control.Block {
			delete(p.ongoing, p.index)
			cancelFunction()
			return nil, errors.New("secuirty blocked")
		}
		sp.Alert = true
	}

	timeoutStr := req.Header.Get("X-Block-Async")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		timeoutStr = "5s"
		timeout, _ = time.ParseDuration(timeoutStr)
	}

	pi.Log.Infof("%s ........... will asynchroniously block after %s", p.name, timeoutStr)

	go func(newCtx context.Context, cancelFunction context.CancelFunc, index uint32, req *http.Request, timeout time.Duration) {
		ticker := time.NewTicker(5 * time.Second)
		sp := p.ongoing[index]
		for {
			select {
			case <-newCtx.Done():
				ticker.Stop()

				if !sp.Alert {
					if sp.Resp.Headers != nil {
						pi.Log.Infof("Done - No Alert! %v", newCtx.Err())
						p.stats("AlertOff")
						if p.wsGate.Control.Learn {
							pi.Log.Infof("index %d\n", index)
							pi.Log.Infof("p.ongoing[index] %v\n", p.ongoing[index])
							pi.Log.Infof("p.pile %v\n", p.pile)
							p.pile.Pile(p.ongoing[index])
							p.reportAllow()
						}
					} else {
						pi.Log.Infof("Done But no Reponse provided - No Alert!")
					}
				} else {
					pi.Log.Infof("Done - With Alert! %v", newCtx.Err())
					p.stats("AlertOn")
				}
				delete(p.ongoing, index)
				return
			case <-time.After(timeout): // maybe reimplement in periodical?
				pi.Log.Infof("Timeout Processing Request!")
				p.stats("Timeout")
				ticker.Stop()
				delete(p.ongoing, index)
				cancelFunction()
				return
			case <-ticker.C:
				if p.periodical(newCtx) {
					pi.Log.Infof("Blocked while processing during tick!")
					if p.wsGate.Control.Block {
						ticker.Stop()
						delete(p.ongoing, p.index)
						cancelFunction()
						return
					}
					sp.Alert = true
				}
			}
		}
	}(newCtx, cancelFunction, p.index, req, timeout)
	p.index++
	return req, nil
}

func (p *plug) ApproveResponse(req *http.Request, resp *http.Response) (*http.Response, error) {
	testBodyHist := true

	//pi.Log.Infof("%s: ApproveResponse started", p.name)

	//pp := new(spec.ProcessProfile)

	now := time.Now()
	ctx := req.Context()
	index, okIndex := ctx.Value(ctxKey("ReqIndex")).(uint32)
	if !okIndex {
		pi.Log.Infof("%s ........... Blocked During Response! Missing context!", p.name)
		return nil, errors.New("missing context")
	}
	sp, spExists := p.ongoing[index]
	if !spExists {
		pi.Log.Infof("%s ........... Blocked During Response! Missing Session!", p.name)
		return nil, errors.New("missing session")
	}

	//reqTime, okReqTime := ctx.Value(ctxKey("ReqTime")).(time.Time)
	//if !okReqTime {
	//	reqTime = now
	//}

	//if req.Header.Get("X-Block-Resp") != "" {
	//	pi.Log.Infof("%s ........... Blocked During Response! returning an error!", p.name)
	//	return nil, errors.New("response blocked")
	//}

	//for name, values := range resp.Header {
	// Loop over all values for the name.
	//	for _, value := range values {
	//	pi.Log.Debugf("%s Response Header: %s: %s", p.name, name, value)
	//	}
	//}

	if p.screenResponse(resp, &sp.Resp) != nil {
		if p.wsGate.Control.Block {
			delete(p.ongoing, p.index)
			return nil, errors.New("secuirty blocked")
		} else {
			sp.Alert = true
		}
	}

	sp.Process.Profile(sp.ReqTime, now, now)
	//fmt.Println(sp.Process.Marshal(0))

	if p.periodical(ctx) {
		pi.Log.Infof("Blocked on periodical during response!")
		if p.wsGate.Control.Block {
			delete(p.ongoing, p.index)
			return nil, errors.New("secuirty blocked")
		}
		sp.Alert = true
	}

	if testBodyHist && resp.Body != nil {
		//fmt.Printf("Analyze Start\n")

		// Asynchrniously stream bytes from the original resp.Body
		// to a new resp.Body
		resp.Body = iofilter.New(resp.Body, responseFilter, sp.RespBody)
	}

	return resp, nil
}

/*
func (p *plug) consultOnRequest(reqProfile *spec.ReqProfile) string {
	postBody, marshalErr := json.Marshal(reqProfile)
	if marshalErr != nil {
		log.Printf("consultOnRequest error during marshal: %v", marshalErr)
		return fmt.Sprintf("Cant marshal in consultOnRequest %v", marshalErr)
	}
	reqBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest(http.MethodPost, p.guardUrl+"/req", reqBody)
	if err != nil {
		pi.Log.Infof("wsgate consultOnRequest: http.NewRequest error %v", err)
	}
	query := req.URL.Query()
	query.Add("sid", p.serviceId)
	query.Add("ns", p.namespace)
	req.URL.RawQuery = query.Encode()

	res, postErr := p.httpc.Do(req)
	//res, postErr := p.httpc.Post(p.guardUrl+"/req", "application/json", reqBody)
	if postErr != nil {
		pi.Log.Infof("wsgate consultOnRequest: httpc.Do error %v", postErr)
		return fmt.Sprintf("Guard unavaliable during consult %v", postErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		pi.Log.Infof("wsgate consultOnRequest: response error %v", readErr)
		return fmt.Sprintf("Guard ilegal response during consult %v", readErr)
	}
	if len(body) != 0 {
		pi.Log.Infof("wsgate consultOnRequest: response is %s", string(body))
		return fmt.Sprintf("Guard: %s", string(body))
	}
	pi.Log.Infof("wsgate consultOnRequest: approved!")
	return ""
}
*/

func (p *plug) reportAllowedPile(pile *spec.Pile) string {
	pi.Log.Infof("reportAllowedPile")
	postBody, marshalErr := json.Marshal(pile)

	if marshalErr != nil {
		pi.Log.Infof("reportAllowedPile error during marshal: %v", marshalErr)
		return fmt.Sprintf("Cant marshal in reportAllowedPile %v", marshalErr)
	}
	reqBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest(http.MethodPost, p.guardUrl+"/pile", reqBody)
	if err != nil {
		pi.Log.Infof("wsgate reportAllowedPile: http.NewRequest error %v", err)
	}
	query := req.URL.Query()
	query.Add("sid", p.serviceId)
	query.Add("ns", p.namespace)
	if p.cmname {
		query.Add("cm", "true")
	}
	req.URL.RawQuery = query.Encode()

	res, postErr := p.httpc.Do(req)
	//res, postErr := p.httpc.Post(p.guardUrl+"/req", "application/json", reqBody)
	if postErr != nil {
		pi.Log.Infof("wsgate reportAllowedPile: httpc.Do error %v", postErr)
		return fmt.Sprintf("Guard unavaliable during consult %v", postErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		pi.Log.Infof("wsgate reportAllowedPile: response error %v", readErr)
		return fmt.Sprintf("Guard ilegal response during consult %v", readErr)
	}
	if len(body) != 0 {
		pi.Log.Infof("wsgate reportAllowedPile: response is %s", string(body))
		return fmt.Sprintf("Guard: %s", string(body))
	}
	pi.Log.Infof("wsgate reportAllowedPile: approved!")

	return ""
}

func (p *plug) stats(key string) {
	p.statistics[key]++
	// build statistics on blocked requests
}

func (p *plug) reportAllow() {
	// send statistics on allowed requests
	p.cycle--
	pi.Log.Infof("reportAllow cycle %d", p.cycle)
	if p.cycle <= 0 {

		p.reportAllowedPile(&p.pile)
		//p.pile.Clear()
		p.cycle = 0
	}
}

func (p *plug) Init() {
	var servicename, namespace, guardUrl string
	pi.Log.Infof("plug %s: Initializing - version %v", p.name, p.version)

	guardUrl = os.Getenv("WSGATE_GUARD_URL")
	if guardUrl == "" {
		// use default
		guardUrl = "http://ws.knative-guard"
	}
	pi.Log.Infof("guardUrl %s", guardUrl)

	namespace = os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = os.Getenv("SERVING_NAMESPACE")
	}
	if namespace == "" {
		// mandatory
		panic("Can't find mandatory parameter namespace")
	}

	servicename = os.Getenv("SERVICENAME")
	if servicename == "" {
		dat, err := os.ReadFile("/etc/podinfo/app")
		if err == nil {
			servicename = string(dat)
		}
	}
	if servicename == "" {
		servicename = os.Getenv("SERVING_SERVICE")
	}
	if servicename == "" {
		// mandatory
		panic("Can't find mandatory parameter  servicename")
	}
	// TBD - change CMNAME to be a boolean as it is derivable from servicename
	// this means servicename should never be "ns.{namespace}"
	cmname := os.Getenv("CMNAME")

	if servicename == "ns."+namespace {
		// mandatory
		panic("Ilegal SERVICENAME - ns.{namespace} is reserved")
	}

	pi.Log.Infof("wsgate configuration: servicename=%s, namespace=%s, cmname=%s, guardUrl=%s", servicename, namespace, cmname, guardUrl)
	p.serviceId = servicename
	p.namespace = namespace
	p.cmname = strings.EqualFold(cmname, "true")
	p.guardUrl = guardUrl

	pi.Log.Infof("wsgate configuration: p.cmname=%t", p.cmname)
	p.kubemgr.InitConfigs()
	p.wsGate = p.kubemgr.FetchConfig(p.namespace, p.serviceId, p.cmname)
}

func init() {
	//fmt.Printf("WSGATE!!!! Initializing!!!!!!!!!<<<<<<<<<<<__________________>>>>>>>>>>\n")
	p := new(plug)
	p.version = version
	p.name = name
	p.pile.Clear()
	p.ongoing = make(map[uint32]*spec.SessionProfile)
	p.statistics = make(map[string]uint32, 8)
	pi.RegisterPlug(p)
	//fmt.Printf("WSGATE!!!! Ended Initializing!!!!!!!!!<<<<<<<<<<<__________________>>>>>>>>>>\n")
}
