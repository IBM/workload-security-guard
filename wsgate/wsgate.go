package wsgate

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	spec "github.com/IBM/workload-security-guard/pkg/apis/wsecurity/v1"
	guardianclientset "github.com/IBM/workload-security-guard/pkg/generated/clientset/guardians"
	v1 "github.com/IBM/workload-security-guard/pkg/generated/clientset/guardians/typed/wsecurity/v1"

	"github.com/IBM/go-security-plugs/iofilter"
	pi "github.com/IBM/go-security-plugs/pluginterfaces"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const version string = "0.0.1"
const name string = "wsgate"

//type StateReq time.Time
//type StateResp time.Time

type ctxKey string

type plug struct {
	name    string
	version string

	// Add here any other state the extension needs
	guardUrl  string
	namespace string
	serviceId string
	//kClient   corev1.ConfigMapInterface
	gClient v1.WsecurityV1Interface
	wsGate  *spec.WsGate
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

func (p *plug) screenRequest(req *http.Request, rp *spec.ReqProfile) error {
	var decission string

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
	fmt.Println(rp.Marshal(0))
	fmt.Println(p.wsGate.Configured)
	ctrl := p.wsGate.Control
	if ctrl.Auto {
		var description bytes.Buffer
		var result string
		for i, leanred := range p.wsGate.Learned {
			if !leanred.Active {
				continue
			}
			result = leanred.Req.Decide(rp)
			if result == "" {
				break
			}
			description.WriteString(fmt.Sprintf("  Learned[%d]: %s", i, result))
		}
		if result != "" {
			// include all decissions!
			decission += description.String()
		}
	} else { // Manual
		if (p.wsGate.Configured != nil) && p.wsGate.Configured.Active {
			decission += p.wsGate.Configured.Req.Decide(rp)
		}

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
	fmt.Println(rp.Marshal(0))
	fmt.Println(p.wsGate.Configured)
	ctrl := p.wsGate.Control
	var decission string
	if ctrl.Auto {
		var description bytes.Buffer
		var result string
		for i, leanred := range p.wsGate.Learned {
			if !leanred.Active {
				continue
			}
			result = leanred.Resp.Decide(rp)
			if result == "" {
				break
			}
			description.WriteString(fmt.Sprintf("  Learned[%d]: %s", i, result))
		}
		if result != "" {
			// include all decissions!
			decission = description.String()
		}
	} else {
		if (p.wsGate.Configured != nil) && p.wsGate.Configured.Active {
			decission = p.wsGate.Configured.Resp.Decide(rp)
		}
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

func responseFilter(buf []byte) error {
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
	fmt.Printf("responseFilter Histogram: %v\n", h)

	return nil
}

func requestFilter(buf []byte) error {
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
	fmt.Printf("requestFilter Histogram: %v\n", h)

	return nil
}

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
	fmt.Printf("sp.ReqTime %v now %v\n", sp.ReqTime, now)
	pp.Profile(sp.ReqTime, sp.ReqTime, now)
	fmt.Println(pp.Marshal(0))

	ctrl := p.wsGate.Control
	var decission string
	if ctrl.Auto {
		var description bytes.Buffer
		var result string
		for i, leanred := range p.wsGate.Learned {
			if !leanred.Active {
				continue
			}
			result = leanred.Process.Decide(pp)
			if decission == "" {
				break
			}
			description.WriteString(fmt.Sprintf("  Learned[%d]: %s", i, result))
		}
		if result != "" {
			// include all decissions!
			decission = description.String()
		}
	} else {
		if (p.wsGate.Configured != nil) && p.wsGate.Configured.Active {
			decission = p.wsGate.Configured.Process.Decide(pp)
		}
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
	testBodyHist := true

	pi.Log.Debugf("%s: ApproveRequest started", p.name)

	//if req.Header.Get("X-Block-Req") != "" {
	//	pi.Log.Infof("%s ........... Blocked During Request! returning an error!", p.name)
	//	return nil, errors.New("request blocked")
	//}

	for name, values := range req.Header {
		// Loop over all values for the name.
		for _, value := range values {
			pi.Log.Debugf("%s Request Header: %s: %s", p.name, name, value)
		}
	}
	p.fetchConfig()
	sp := new(spec.SessionProfile)
	sp.ReqTime = time.Now()
	if p.screenRequest(req, &sp.Req) != nil {
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

	if testBodyHist && req.Body != nil {
		//fmt.Printf("Analyze Start\n")

		// Asynchrniously stream bytes from the original resp.Body
		// to a new resp.Body
		req.Body = iofilter.New(req.Body, requestFilter)
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
					pi.Log.Infof("Done - No Alert! %v", newCtx.Err())
					p.stats("AlertOff")
					if p.wsGate.Control.Learn {
						p.pile.Pile(p.ongoing[index])
						p.reportAllow()
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

	pi.Log.Infof("%s: ApproveResponse started", p.name)

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

	for name, values := range resp.Header {
		// Loop over all values for the name.
		for _, value := range values {
			pi.Log.Infof("%s Response Header: %s: %s", p.name, name, value)
		}
	}

	if p.screenResponse(resp, &sp.Resp) != nil {
		if p.wsGate.Control.Block {
			delete(p.ongoing, p.index)
			return nil, errors.New("secuirty blocked")
		} else {
			sp.Alert = true
		}
	}

	sp.Process.Profile(sp.ReqTime, now, now)
	fmt.Println(sp.Process.Marshal(0))

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
		resp.Body = iofilter.New(resp.Body, responseFilter)
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
	postBody, marshalErr := json.Marshal(pile)
	if marshalErr != nil {
		log.Printf("reportAllowedPile error during marshal: %v", marshalErr)
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
	if p.cycle <= 0 {
		p.reportAllowedPile(&p.pile)
		//p.pile.Clear()
		p.cycle = 0
	}
}

func (p *plug) initCrd() {
	var kubeconfig *string
	var cfg *rest.Config
	var errInCluster error
	var errOutOfCluster error
	// creates the in-cluster config

	cfg, errInCluster = rest.InClusterConfig()
	if errInCluster != nil {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()
		// use the current context in kubeconfig
		cfg, errOutOfCluster = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if errOutOfCluster != nil {
			panic(fmt.Sprintf("No Config found! errInCluster %s", errInCluster.Error()))
		}
	}

	guardianClient, err := guardianclientset.NewForConfig(cfg)
	if err != nil {
		panic(err.Error())
	}

	p.gClient = guardianClient.WsecurityV1()

	//clientset, err := kubernetes.NewForConfig(cfg)
	//if err != nil {
	//	panic(err.Error())
	//}
	//p.kClient = clientset.CoreV1().ConfigMaps("knative-serving")

}

//
func (p *plug) readCrd(namespace string, serviceId string) *spec.GuardianSpec {
	g, err := p.gClient.Guardians(namespace).Get(context.TODO(), serviceId, metav1.GetOptions{})
	if err != nil {
		pi.Log.Infof("Err during get %s.%s: %s", serviceId, namespace, err.Error())
		//panic(fmt.Sprintf("No Guardian! for %s.%s", serviceId, namespace))
		return nil
	}
	pi.Log.Infof("Found guardian %s.%s", serviceId, namespace)
	fmt.Print((*spec.WsGate)(g.Spec).Marshal(0))
	(*spec.WsGate)(g.Spec).Reconcile()
	return g.Spec
}

/*func (p *plug) readConfigMap() {
	cm, err := p.kClient.Get(context.TODO(), "guardian", metav1.GetOptions{})
	if err != nil {
		fmt.Printf("ConfigMap Error: %v\n", err)
		panic(err.Error())
	}
	p.blockByDefault = cm.Data["BlockByDefault"] != "false"
	fmt.Printf("ConfigMap: %s\n", cm.Data["BlockByDefault"])

	//	err = json.Unmarshal([]byte(cm.Data["guardian"]), &p.wsGate)
	//	if err != nil {
	//		fmt.Printf("ConfigMap Unmarshal Error: %v\n", err)
	//		panic(err.Error())
	//	}
}
*/
func (p *plug) fetchConfig() {
	//p.readConfigMap()
	gurdianSpec := p.readCrd(p.namespace, p.serviceId)
	if gurdianSpec == nil {
		gurdianSpec = p.readCrd("knative-serving", "guardian")
	}
	if gurdianSpec == nil {
		fmt.Println("Guardian was not set!")
		gurdianSpec = new(spec.GuardianSpec)
		// default gurdianSpec has:
		// 		gurdianSpec.falseAllow=false
		// 		gurdianSpec.ConsultGuard.Active = false
	}
	//if gurdianSpec.Configured == nil {
	//	gurdianSpec.Configured = new(spec.Critiria)
	//}
	//if gurdianSpec.Learned == nil {
	//	gurdianSpec.Learned = make([]spec.Critiria, 0)
	//}
	p.wsGate = (*spec.WsGate)(gurdianSpec)
	fmt.Printf("p.wsGate %v", p.wsGate)

	/*
		req, err := http.NewRequest(http.MethodGet, p.guardUrl+"/config", nil)
		if err != nil {
			pi.Log.Infof("wsgate getConfig: http.NewRequest error %v", err)
		}
		query := req.URL.Query()
		query.Add("sid", p.id)
		req.URL.RawQuery = query.Encode()
		res, getErr := p.httpc.Do(req)
		if getErr != nil {
			pi.Log.Infof("wsgate getConfig: httpc.Do error %v", getErr)
			return
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			pi.Log.Infof("wsgate getConfig: read body error %v", readErr)
		}

		//pi.Log.Infof("wsgate getConfig: unmarshal %s", string(body))
		jsonErr := json.Unmarshal(body, &p.gateConfig)
		if jsonErr != nil {
			pi.Log.Infof("wsgate getConfig: unmarshel error %v", jsonErr)
		}

		//pi.Log.Infof("wsgate getConfig: ended %v ", p.gateConfig)
	*/
}

func (p *plug) Init() {
	pi.Log.Infof("plug %s: Initializing - version %v", p.name, p.version)

	p.guardUrl = os.Getenv("WSGATE_GUARD_URL")
	if p.guardUrl == "" {
		p.guardUrl = "http://ws.knative-guard"
	}
	pi.Log.Infof("guardUrl %s", p.guardUrl)

	servingNamespace := os.Getenv("SERVING_NAMESPACE")
	if servingNamespace == "" {
		panic("Cant find SERVING_NAMESPACE")
	}
	servingService := os.Getenv("SERVING_SERVICE")
	if servingService == "" {
		panic("Cant find SERVING_SERVICE")
	}
	p.serviceId = servingService
	p.namespace = servingNamespace

	p.initCrd()
	p.fetchConfig()
}

func init() {
	fmt.Printf("WSGATE!!!! Initializing!!!!!!!!!<<<<<<<<<<<__________________>>>>>>>>>>\n")
	p := new(plug)
	p.version = version
	p.name = name
	p.pile.Clear()
	p.ongoing = make(map[uint32]*spec.SessionProfile)
	p.statistics = make(map[string]uint32, 8)
	pi.RegisterPlug(p)
	fmt.Printf("WSGATE!!!! Ended Initializing!!!!!!!!!<<<<<<<<<<<__________________>>>>>>>>>>\n")
}
