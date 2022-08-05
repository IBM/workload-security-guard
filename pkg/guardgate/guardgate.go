package guardgate

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
	"runtime/debug"
	"strings"
	"time"

	spec "github.com/IBM/workload-security-guard/pkg/apis/wsecurity/v1"
	"github.com/IBM/workload-security-guard/pkg/guardkubemgr"

	"github.com/IBM/go-security-plugs/iodup"
	"github.com/IBM/go-security-plugs/iofilter"
	pi "github.com/IBM/go-security-plugs/pluginterfaces"

	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

const plugVersion string = "0.0.1"
const plugName string = "guard"

const (
	GuardianLoadIntervalDefault = 5 * time.Minute
	ReportPileIntervalDefault   = 10 * time.Minute
	PodMonitorIntervalDefault   = 5 * time.Second
	MinimumInterval             = 5 * time.Second
)

var errSecurity error = errors.New("secuirty blocked by guard")

type ctxKey string

type plug struct {
	name        string
	version     string
	serviceName string
	namespace   string
	config      map[string]string
	ctx         context.Context

	kubemgr guardkubemgr.Kubemgr

	// Add here any other state the extension needs
	guardUrl             string
	useConfigmap         bool
	monitorPod           bool
	wsGate               *spec.GuardianSpec
	criteria             *spec.Criteria
	GuardianLoadInterval time.Duration
	ReportPileInterval   time.Duration
	PodMonitorInterval   time.Duration
	httpc                http.Client
	pile                 spec.Pile
	pileCount            int
	statistics           map[string]uint32
	pod                  spec.PodProfile
	alert                string
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))[0:16]
}

func (p *plug) Shutdown() {
	pi.Log.Infof("%s: Shutdown", p.name)
	p.reportPile()
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

func (p *plug) screenResponseBody(req *http.Response, bp *spec.BodyProfile) string {
	p.stats("RespBodyOk")
	return ""
}
func (p *plug) screenRequestBody(req *http.Request, bp *spec.BodyProfile) string {
	testBodyHist := true
	maxBody := int64(1000000)
	doneBody := false

	if req.ContentLength < maxBody {
		// TBD - validate content-type params returned by ParseMediaType!
		ctype, _, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
		if err != nil {
			ctype = "application/octet-stream"
		}
		method := req.Method
		//pi.Log.Debugf("%s Request Content-Type: %s  (params %v)\n", p.name, ctype, params)
		//pi.Log.Debugf("%s Request Method: %s\n", p.name, method)
		if ctype == "application/json" {
			//pi.Log.Debugf("Processing json!\n")
			doneBody = true
			dup := iodup.New(req.Body, 2, 128, 8192)

			var j interface{}
			dec := json.NewDecoder(&dup.Output[1])
			err = dec.Decode(&j)
			if err != nil {
				pi.Log.Debugf("Error Decoding json! %v\n", err)
			}
			bp.Structured = new(spec.StructuredProfile)
			bp.Structured.Profile(j)
			req.Body = &dup.Output[0]

		} else if method == "POST" || method == "PUT" || method == "PATCH" {
			//pi.Log.Debugf("Processing POST/PUT/PATCH\n")
			if ctype == "application/x-www-form-urlencoded" {
				//pi.Log.Debugf("Processing application/x-www-form-urlencoded!\n")
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
				//pi.Log.Debugf("Form Data %s\n", req.Form)

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
			} else if ctype == "multipart/form-data" {
				//pi.Log.Debugf("Processing multipart/form-data!\n")
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
				//pi.Log.Debugf("Form Data %s\n", req.Form)

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
			pi.Log.Debugf("Analyzing Body")

			// Asynchrniously stream bytes from the original resp.Body
			// to a new resp.Body
			//req.Body = iofilter.New(req.Body, requestFilter, sp.ReqBody)
		}
	}

	if (p.criteria != nil) && p.criteria.Active {
		if decission := p.criteria.ReqBody.Decide(bp); decission != "" {
			p.stats("ReqBodyNok")
			return fmt.Sprintf("HttpRequestBody: %s", decission)
		}
	}
	p.stats("ReqBodyOk")
	return ""
}

func (p *plug) screenEnvelop(sp *spec.SessionProfile) string {
	now := time.Now()
	sp.Envelop.Profile(sp.ReqTime, now, now)

	if (p.criteria != nil) && p.criteria.Active {
		if decission := p.criteria.Envelop.Decide(&sp.Envelop); decission != "" {
			p.stats("EnvelopNok")
			return fmt.Sprintf("Envelop: %s", decission)
		}
	}
	p.stats("EnvelopOk")
	return ""
}

func (p *plug) screenRequest(req *http.Request, rp *spec.ReqProfile) string {
	var decission string

	//pi.Log.Debugf("screenRequest: ctrl %v", p.wsGate.Control)
	p.stats("Total")
	// Request client and server identities
	cip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		decission += fmt.Sprintf("illegal req.RemoteAddr %s", err.Error())
	}
	_, _, err = net.SplitHostPort(req.URL.Host)

	if err != nil {
		decission += fmt.Sprintf("illegal req.URL.Host %s", err.Error())
	}
	//pi.Log.Debugf("Client: %s port %s", cip, cport)
	//pi.Log.Debugf("Server: %s port %s", sip, sport)

	// Request principles
	//pi.Log.Debugf("req.Method %s", req.Method)
	//pi.Log.Debugf("req.Proto %s", req.Proto)
	//pi.Log.Debugf("scheme: %s", req.URL.Scheme)
	//pi.Log.Debugf("opaque: %s", req.URL.Opaque)

	//pi.Log.Debugf("ContentLength: %d", req.ContentLength)
	//pi.Log.Debugf("Trailer: %#v", req.Trailer)

	// TBD req.Form

	ip := net.ParseIP(cip)
	rp.Profile(req, ip)

	if (p.criteria != nil) && p.criteria.Active {
		if decission := p.criteria.Req.Decide(rp); decission != "" {
			p.stats("ReqNok")
			return fmt.Sprintf("HttpRequest: %s", decission)
		}
	}

	p.stats("ReqOk")
	return ""
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

func (p *plug) screenResponse(resp *http.Response, rp *spec.RespProfile) string {
	rp.Profile(resp)

	if (p.criteria != nil) && p.criteria.Active {
		if decission := p.criteria.Resp.Decide(rp); decission != "" {
			p.stats("RespNok")
			return fmt.Sprintf("HttpResponse: %s", decission)
		}
	}

	p.stats("RespOk")
	return ""
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

func (p *plug) ApproveRequest(req *http.Request) (*http.Request, error) {
	//pi.Log.Debugf("%s: ApproveRequest started", p.name)

	sp := new(spec.SessionProfile)
	sp.ReqTime = time.Now()

	// Req
	sp.Alert += p.screenEnvelop(sp)
	sp.Alert += p.screenRequest(req, &sp.Req)
	sp.Alert += p.screenRequestBody(req, &sp.ReqBody)
	if (sp.Alert != "" || p.alert != "") && p.wsGate.Control.Block {
		p.stats("BlockOnRequest")
		sp.Cancel()
		return nil, errSecurity
	}

	ctx, cancelFunction := context.WithCancel(req.Context())
	sp.Cancel = cancelFunction

	//ctx = context.WithValue(ctx, ctxKey("ReqIndex"), index)
	ctx = context.WithValue(ctx, ctxKey("GuardSession"), sp)
	req = req.WithContext(ctx)

	//goroutine to accompany the request
	go func(ctx context.Context, sp *spec.SessionProfile) {
		sp.SecurityMonitorTicker = time.NewTicker(p.PodMonitorInterval)
		done := false
		for !done {
			done = p.securityMonitor(ctx, sp)
			if !done {
				sp.SecurityMonitorTicker.Stop()
				pi.Log.Debugf("Making a second attampt at securityMonitor")
				p.securityMonitor(ctx, sp)
				done = true
			}
		}
	}(ctx, sp)
	return req, nil
}

func (p *plug) securityMonitor(ctx context.Context, sp *spec.SessionProfile) bool {
	defer func() {
		if r := recover(); r != nil {
			pi.Log.Warnf("securityMonitor Recovered %s", r)
			pi.Log.Debugf(string(debug.Stack()))
		}
	}()
	for {
		select {
		case <-ctx.Done(): // Always finish any request here!
			sp.SecurityMonitorTicker.Stop()

			// Should we learn?
			if p.wsGate.Control.Learn && (sp.Alert == "" || p.wsGate.Control.Force) && sp.Resp.Headers != nil {
				// Learn only if asked to learn and we received a response
				if p.monitorPod {
					p.pile.Pile(sp, &p.pod)
				} else {
					p.pile.Pile(sp, nil)
				}

				p.pileCount += 1
				pi.Log.Debugf("Learn - add to pile! pileCount %d", p.pileCount)
			}

			// Should we alert?
			if sp.Alert != "" {
				p.Alert(sp.Alert)
				p.stats("SessionLevelAlert")
			} else {
				if sp.Resp.Headers != nil {
					pi.Log.Debugf("No Alert!")
					p.stats("NoAlert")
				} else {
					pi.Log.Debugf("No Alert while completed before receiving a response!")
					p.stats("NoResponse")
				}
			}
			return true
		case <-sp.SecurityMonitorTicker.C:
			sp.Alert += p.screenEnvelop(sp)
			if p.wsGate.Control.Block {
				if sp.Alert != "" {
					sp.Cancel()
				}
				if p.alert != "" {
					p.stats("BlockOnPod")
					sp.Cancel()
				}
			}

		}
	}
}

func (p *plug) Alert(alert string) {
	pi.Log.Infof("SECURITY ALERT! %s", alert)
}

func (p *plug) ApproveResponse(req *http.Request, resp *http.Response) (*http.Response, error) {
	testBodyHist := true

	//pi.Log.Infof("%s: ApproveResponse started", p.name)

	//pp := new(spec.ProcessProfile)

	ctx := req.Context()
	/*
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
	*/
	sp, spExists := ctx.Value(ctxKey("GuardSession")).(*spec.SessionProfile)
	if !spExists { // This should never happen!
		pi.Log.Warnf("%s ........... Blocked During Response! Missing context!", p.name)
		return nil, errors.New("missing context")
	}

	sp.Alert += p.screenEnvelop(sp)
	sp.Alert += p.screenResponse(resp, &sp.Resp)
	sp.Alert += p.screenResponseBody(resp, &sp.RespBody)
	if (sp.Alert != "" || p.alert != "") && p.wsGate.Control.Block {
		sp.Cancel()
		p.stats("BlockOnResponse")
		return nil, errSecurity
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
		pi.Log.Infof("guardgate consultOnRequest: http.NewRequest error %v", err)
	}
	query := req.URL.Query()
	query.Add("sid", p.serviceId)
	query.Add("ns", p.namespace)
	req.URL.RawQuery = query.Encode()

	res, postErr := p.httpc.Do(req)
	//res, postErr := p.httpc.Post(p.guardUrl+"/req", "application/json", reqBody)
	if postErr != nil {
		pi.Log.Infof("guardgate consultOnRequest: httpc.Do error %v", postErr)
		return fmt.Sprintf("Guard unavaliable during consult %v", postErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		pi.Log.Infof("guardgate consultOnRequest: response error %v", readErr)
		return fmt.Sprintf("Guard ilegal response during consult %v", readErr)
	}
	if len(body) != 0 {
		pi.Log.Infof("guardgate consultOnRequest: response is %s", string(body))
		return fmt.Sprintf("Guard: %s", string(body))
	}
	pi.Log.Infof("guardgate consultOnRequest: approved!")
	return ""
}
*/

func (p *plug) clearPile() {
	p.pile.Clear()
	p.pileCount = 0
}

func (p *plug) reportPile() {
	if p.pileCount == 0 {
		pi.Log.Debugf("No pile to report to guard-service!")
		return
	}
	defer p.clearPile()

	pi.Log.Infof("Reporting a pile with pileCount %d records to guard-service", p.pileCount)

	postBody, marshalErr := json.Marshal(p.pile)

	if marshalErr != nil {
		pi.Log.Warnf("Error during marshal: %v", marshalErr)
		return
	}
	reqBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest(http.MethodPost, p.guardUrl+"/pile", reqBody)
	if err != nil {
		pi.Log.Warnf("Http.NewRequest error %v", err)
		return
	}
	query := req.URL.Query()
	query.Add("sid", p.serviceName)
	query.Add("ns", p.namespace)
	if p.useConfigmap {
		query.Add("cm", "true")
	}
	req.URL.RawQuery = query.Encode()

	res, postErr := p.httpc.Do(req)
	if postErr != nil {
		pi.Log.Warnf("Httpc.Do error %v", postErr)
		return
	}

	defer res.Body.Close()
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		pi.Log.Infof("Response error %v", readErr)
		return
	}
	if len(body) != 0 {
		pi.Log.Infof("guard-service response is %s", string(body))
	}

}

func (p *plug) stats(key string) {
	p.statistics[key]++
	// build statistics on blocked requests
}

func (p *plug) Start(ctx context.Context) context.Context {
	p.ctx = ctx
	return ctx
}

func (p *plug) loadGuardian() {
	pi.Log.Infof("(Re)loading Guardian")
	p.wsGate = p.kubemgr.FetchConfig(p.namespace, p.serviceName, p.useConfigmap)
	if p.wsGate.Control.Alert {
		if p.wsGate.Control.Auto {
			p.criteria = p.wsGate.Learned
		} else {
			p.criteria = p.wsGate.Configured
		}
	}
}

func (p *plug) decidePod() {
	if !p.monitorPod {
		return
	}
	p.alert = ""
	p.pod.Profile()
	if p.criteria != nil {
		if decission := p.criteria.Pod.Decide(&p.pod); decission != "" {
			p.alert = fmt.Sprintf("Pod: %s", decission)
			p.Alert(p.alert)
			p.stats("PodLevelAlert")
		}
	}
}

func parseInterval(name string, intervalStr string, defaultVal time.Duration) (d time.Duration) {
	var err error

	d = defaultVal
	if intervalStr != "" {
		d, err = time.ParseDuration(intervalStr)
	}
	if err != nil {
		pi.Log.Infof("Interval %s ilegal value %s - using defualt value instead (Err: %v)", name, intervalStr, err)
		d = defaultVal
	}
	if d < MinimumInterval {
		pi.Log.Infof("Interval %s value %s is too low, using minimum value instead", name, intervalStr)
		d = MinimumInterval
	}
	return
}

func (p *plug) Init(ctx context.Context, c map[string]string, serviceName string, namespace string, logger pi.Logger) context.Context {
	var ok bool
	var v string

	p.ctx = ctx
	p.serviceName = serviceName
	p.namespace = namespace
	p.config = c

	if p.guardUrl, ok = c["guard-url"]; !ok {
		// use default
		p.guardUrl = "http://guard-service.knative-guard"
	}

	v, ok = c["use-configmap"]
	if ok && strings.EqualFold(v, "true") {
		p.useConfigmap = true
	}

	v, ok = c["monitor-pod"]
	if ok && strings.EqualFold(v, "true") {
		p.monitorPod = true
	}

	// svcname should never be "ns.{namespace}" as this is a reserved name
	if p.serviceName == "ns."+p.namespace {
		// mandatory
		panic("Ilegal Svcname - ns.{Namespace} is reserved")
	}

	pi.Log.Debugf("guardgate configuration: servicename=%s, namespace=%s, cmname=%t, guardUrl=%s, p.monitorPod=%t, guardian-load-interval %v, report-pile-interval %v, pod-monitor-interval %v",
		p.serviceName, p.namespace, p.useConfigmap, p.guardUrl, p.monitorPod, c["guardian-load-interval"], c["report-pile-interval"], c["pod-monitor-interval"])

	p.clearPile()
	p.statistics = make(map[string]uint32, 8)

	p.kubemgr.InitConfigs()

	p.GuardianLoadInterval = parseInterval("GuardianLoad", c["guardian-load-interval"], GuardianLoadIntervalDefault)
	p.ReportPileInterval = parseInterval("ReportPile", c["report-pile-interval"], ReportPileIntervalDefault)
	p.PodMonitorInterval = parseInterval("PodMonitor", c["pod-monitor-interval"], PodMonitorIntervalDefault)

	p.loadGuardian()
	p.decidePod()

	//goroutine for Guard instance
	go func() {
		var guardianLoadTicker, reportPileTicker, podMonitorTicker *time.Ticker
		if p.GuardianLoadInterval != time.Duration(0) {
			guardianLoadTicker = time.NewTicker(p.GuardianLoadInterval)
		}
		if p.ReportPileInterval != time.Duration(0) {
			reportPileTicker = time.NewTicker(p.ReportPileInterval)
		}
		if p.PodMonitorInterval != time.Duration(0) {
			podMonitorTicker = time.NewTicker(p.PodMonitorInterval)
		}
		for {
			select {
			case <-p.ctx.Done(): // Always finish guard here!
				pi.Log.Infof("%s Done!", plugName)
				if guardianLoadTicker != nil {
					guardianLoadTicker.Stop()
				}
				if reportPileTicker != nil {
					reportPileTicker.Stop()
				}
				if podMonitorTicker != nil {
					podMonitorTicker.Stop()
				}
				p.reportPile()
				pi.Log.Debugf("Statistics %v", p.statistics)
				return
			case <-guardianLoadTicker.C:
				pi.Log.Debugf("Load Guardian Ticker")
				p.loadGuardian()
				pi.Log.Debugf("Statistics %v", p.statistics)
			case <-reportPileTicker.C:
				pi.Log.Debugf("Report Pile Ticker")
				p.reportPile()
			case <-podMonitorTicker.C:
				p.decidePod()
			}
		}
	}()
	return ctx
}

func init() {
	pi.RegisterPlug(&plug{version: plugVersion, name: plugName})
}
