package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	spec "github.com/IBM/workload-security-guard/pkg/apis/wsecurity/v1"
	"github.com/IBM/workload-security-guard/pkg/guardkubemgr"

	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

type serviceRecord struct {
	wsgate *spec.GuardianSpec
	pile   spec.Pile
}

type learner struct {
	services   map[string]*serviceRecord
	kmgr       guardkubemgr.Kubemgr
	namespaces map[string]bool
}

//var services = make(map[string]*serviceRecord)
//var kmgr wsgate.Kubemgr

func processQuery(query url.Values) (ns string, sid string, cmflag bool) {
	fmt.Printf("Servicing processQuery %v\n", query)

	cmflagSlice := query["cm"]
	sidSlice := query["sid"]
	nsSlice := query["ns"]
	fmt.Printf("Servicing processQuery sidSlice %v nsSlice %v cmflagSlice %v\n", sidSlice, nsSlice, cmflagSlice)
	if len(sidSlice) != 1 || len(nsSlice) != 1 || len(cmflagSlice) > 1 {
		fmt.Printf("Servicing processQuery wrong data sid %d ns %d cmflag %d\n", len(sidSlice), len(nsSlice), len(cmflagSlice))
		return
	}
	sid = sidSlice[0]
	ns = nsSlice[0]
	if len(cmflagSlice) > 0 {
		cmflag = (cmflagSlice[0] == "true")
	}
	if sid == "ns-"+ns {
		fmt.Printf("Servicing processQuery ilegal sid\n")
		sid = ""
		return
	}
	fmt.Printf("Servicing processQuery sid %s ns %s cmflag %t\n", sid, ns, cmflag)
	return
}

func (l *learner) processPile(w http.ResponseWriter, req *http.Request) {
	// Add security to ensure that only gate can use this interface!
	// Check that the source is from the local 10.*.*.* range
	// Check that the
	var pile spec.Pile
	ns, sid, cmflag := processQuery(req.URL.Query())
	if sid == "" || ns == "" {
		http.Error(w, "processPile Missing data", http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(req.Body).Decode(&pile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//fmt.Printf("Servicing processPile of service id %s:%s (cmflag %t) %v\n", ns, sid, cmflag, pile)
	record := l.loadSession(ns, sid, cmflag)
	//fmt.Printf("Servicing processPile loadSession %v\n", record)

	record.pile.Append(&pile)
	//fmt.Printf("Servicing processPile Record after append %v\n", record.pile)
	//fmt.Printf("Servicing processPile Marshal Record\n______________\n%s\n", record.pile.Marshal())

	learn := new(spec.Criteria)
	//fmt.Printf("Servicing processPile after learn new\n")

	learn.Learn(&record.pile)
	//fmt.Printf("Servicing processPile after learn %v\n", learn.Process.Tcp4Peers)
	if record.wsgate.Learned != nil {
		//fmt.Printf("Servicing processPile before merging %v and %v\n", learn.Process.Tcp4Peers, record.wsgate.Learned.Process.Tcp4Peers)
		learn.Merge(record.wsgate.Learned)
		//fmt.Printf("Servicing processPile after merging %v and %v\n", learn.Process.Tcp4Peers, record.wsgate.Learned.Process.Tcp4Peers)
	}
	record.wsgate.Learned = learn
	record.wsgate.Learned.Active = true
	fmt.Printf("About to setCrd record.wsgate.Control %v\n", record.wsgate.Control)

	//fmt.Printf("Marshal Learned\n______________\n%s\n", record.wsgate.Learned.Marshal(0))
	//fmt.Printf("Marshal Learned: Method %d %v\n", len(record.wsgate.Learned.Req.Method), record.wsgate.Learned.Req.Method)
	if cmflag {
		ret := l.kmgr.SetCm(ns, sid, record.wsgate)
		fmt.Printf("Pile SetCm returned %s\n", ret)
		//updateCm(cmflag, ns, sid, record.wsgate)
	} else {
		ret := l.kmgr.SetCrd(ns, sid, record.wsgate)
		fmt.Printf("Pile setCrd returned %s\n", ret)
		//updateCrd(ns, sid, record.wxsgate)
	}

	data := ""
	w.Write([]byte(data))
}

func (l *learner) consultOnReq(w http.ResponseWriter, req *http.Request) {
	// Add security to ensure that only gate can use this interface!
	// Check that the source is from the local 10.*.*.* range
	// Check that the
	query := req.URL.Query()
	sidSlice := query["sid"]
	nsSlice := query["ns"]
	//fmt.Printf("Servicing consultOnReq %v\n", query)
	if len(sidSlice) != 1 || len(nsSlice) != 1 {
		fmt.Printf("Servicing consultOnReq missing data sid %d ns %d\n", len(sidSlice), len(nsSlice))
		return
	}
	sid := sidSlice[0]
	ns := nsSlice[0]
	if sid == "" || ns == "" {
		fmt.Printf("Servicing consultOnReq missing data\n")
		return
	}
	//fmt.Printf("Servicing consultOnReq of service id %s:%s\n", ns, sid)
	data := "No Way!"

	w.Write([]byte(data))
}

/*
func _setCrd() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	//kubeconfig := "/Users/chris/.kube/config"
	//config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	//if err != nil {
	//	panic(err.Error())
	//}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	data := new(spec.WsGate)
	data.Control.Consult = true
	data.Control.RequestsPerMinuete = 60
	//data.ForceAllow = true
	data.Configured = new(spec.Criteria)
	data.Configured.Req.AddTypicalVal()

	json, err := json.Marshal(data)
	if err != nil {
		panic(err.Error())
	}

	b, err := clientset.RESTClient().Post().AbsPath("apis/cminion.com/v1/nodescost").Body(json).DoRaw(context.TODO())
	if err != nil {
		fmt.Println(string(b))
		panic(err.Error())
	}
	fmt.Printf("b is %s\n", string(b))

}
*/
//var gClient v1.WsecurityV1Interface

/*

func getGuardianClient() {
	kmgr.InitConfigs()
		var kubeconfig *string
		var cfg *rest.Config
		var err error
		cfg, err = rest.InClusterConfig()
		if err != nil {
			if home := homedir.HomeDir(); home != "" {
				kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
			} else {
				kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
			}
			flag.Parse()

			// use the current context in kubeconfig
			cfg, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
			if err != nil {
				panic(err.Error())
			}
		}

		guardianClient, err := guardianclientset.NewForConfig(cfg)
		if err != nil {
			panic(err.Error())
		}

		gClient = guardianClient.WsecurityV1()

}
*/
func (l *learner) storeSession(ns string, sid string, wsgate *spec.GuardianSpec) {
	fmt.Printf("storeSession %s.%s\n", ns, sid)
	if wsgate == nil {
		return
	}
	service := sid + "." + ns
	record, exists := l.services[service]
	if !exists {
		record = new(serviceRecord)
		record.pile.Clear()
		l.services[service] = record
	}
	record.wsgate = wsgate
	if _, ok := l.namespaces[ns]; !ok {
		l.namespaces[ns] = true
		//go l.watchNamespace(ns)
	}
}

func (l *learner) loadSession(ns string, sid string, cmname bool) *serviceRecord {
	service := sid + "." + ns
	record, exists := l.services[service]
	if exists {
		return record
	}

	// not cached
	gate := l.kmgr.FetchConfig(ns, sid, cmname)
	l.storeSession(ns, sid, gate)
	return l.services[service]
}

func (l *learner) deleteSession(ns string, sid string) {
	fmt.Printf("deleteSession %s.%s\n", ns, sid)
	delete(l.services, sid+"."+ns)
}

func (l *learner) set(ns string, sid string, g *spec.GuardianSpec) {
	if g == nil {
		l.deleteSession(ns, sid)
	} else {
		l.storeSession(ns, sid, g)
	}
}

func (l *learner) watchNamespace(ns string) {
	for {
		l.kmgr.WatchOnce(ns, l.set)
		//watchCrdOnce(ns)
		sleep("100s")
	}
}

func sleep(timeoutStr string) {
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		timeoutStr = "1s"
		timeout, _ = time.ParseDuration(timeoutStr)
	}
	time.Sleep(timeout)
}

/*
func watchCrdOnce(ns string) {
	defer func() {
		if recovered := recover(); recovered != nil {
			fmt.Printf("Recovered from panic during watchCrdOnce! Recover: %v\n", recovered)
		}
	}()
	watcher, err := gClient.Guardians(ns).Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("watchCrd gClient.Guardians(%s).Watch err %v\n", ns, err)
		return
	}
	ch := watcher.ResultChan()
	for {
		select {
		case event, ok := <-ch:
			if !ok {
				// the channel got closed, so we need to restart
				fmt.Printf("Kubernetes hung up on us, restarting event watcher\n")
				return
			}

			// handle the event
			fmt.Printf("event %v\n", event)
			//fmt.Printf("--------------> Event\n")

			switch event.Type {
			case watch.Added:
				g := event.Object.(*spec.Guardian)
				ns := g.ObjectMeta.Namespace
				sid := g.ObjectMeta.Name
				fmt.Printf("Service %s/%s added\n", ns, sid)
				storeSession(ns, sid, g.Spec)
			case watch.Modified:
				g := event.Object.(*spec.Guardian)
				ns := g.ObjectMeta.Namespace
				sid := g.ObjectMeta.Name
				fmt.Printf("Service %s/%s modified\n", ns, sid)
				storeSession(ns, sid, g.Spec)
			case watch.Deleted:
				g := event.Object.(*spec.Guardian)
				ns := g.ObjectMeta.Namespace
				sid := g.ObjectMeta.Name
				fmt.Printf("Service %s/%s deleted\n", ns, sid)
				delete(services, sid+"."+ns)
			case watch.Error:
				s := event.Object.(*metav1.Status)
				fmt.Printf("Error during watch: \n\tListMeta %v\n\tTypeMeta %v\n", s.ListMeta, s.TypeMeta)
			}
		case <-time.After(10 * time.Minute):
			// deal with the issue where we get no events
			fmt.Printf("Timeout, restarting event watcher\n")
			return
		}
	}
}
*/
/*
func getCrd(ns string, sid string) *spec.WsGate {
	var g *spec.Guardian
	var err error

	fmt.Printf("getCrd guardian %s %s\n", ns, sid)

	g, err = gClient.Guardians(ns).Get(context.TODO(), sid, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("getCrd err %v\n", err)
		return nil
	}
	return (*spec.WsGate)(g.Spec)
}


func updateCm(cmname string, ns string, sid string, wsGate *spec.WsGate) {

}

func updateCrd(ns string, sid string, wsGate *spec.WsGate) {
	var w spec.WsGate
	w.Learned = wsGate.Learned
	var err error

	fmt.Printf("updateCrd guardian %s %s\n", ns, sid)
	//g.Name = sid
	//g.Spec = (*spec.GuardianSpec)(wsGate)
	//_, err = gClient.Guardians(ns).Update(context.TODO(), &g, metav1.UpdateOptions{})
	buf, err := json.Marshal(w)
	if err != nil {
		fmt.Printf("updateCrd  json.Marshal(wsGate) err %v\n", err)
		return
	}
	_, err = gClient.Guardians(ns).Patch(context.TODO(), sid, types.MergePatchType, buf, metav1.PatchOptions{})
	if err != nil {
		fmt.Printf("updateCrd err %v\n", err)
	}
}

func setCrd(ns string, sid string, wsGate *spec.WsGate) {
	var g *spec.Guardian
	var err error

	fmt.Printf("setCrd guardian %s %s\n", ns, sid)
	g, err = gClient.Guardians(ns).Get(context.TODO(), sid, metav1.GetOptions{})
	if err == nil {
		fmt.Printf("setCrd guardian %v\n", g)
		g.Name = sid
		g.Spec = (*spec.GuardianSpec)(wsGate)
		_, err = gClient.Guardians(ns).Update(context.TODO(), g, metav1.UpdateOptions{})
		if err != nil {
			fmt.Printf("setCrd update err %v\n", err)
		}
	} else {
		fmt.Printf("setCrd get err %v\n", err)
		g.Name = sid
		g.Spec = (*spec.GuardianSpec)(wsGate)
		_, err = gClient.Guardians(ns).Create(context.TODO(), g, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("setCrd create err %v\n", err)
		}
	}
}
*/
/*
func setConfigMap() {
	data := new(spec.WsGate)
	data.Control.Consult = true
	data.Control.RequestsPerMinuete = 60
	//data.ForceAllow = true
	data.Configured = new(spec.Criteria)
	data.Configured.Req.AddTypicalVal()

	databuf, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err.Error())
	}

	configMap := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "guardian",
			Namespace: "default",
		},
		Data: map[string]string{"guardian": string(databuf)},
	}

	fmt.Printf("MarshalIdent\n %s\n", string(databuf))
	fmt.Printf("My Marshal\n %s\n", data.Configured.Req.Marshal(4))

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	//kubeClient, err := apiextension.NewForConfig(config)
	//if err != nil {
	//	log.Fatalf("Failed to create client: %v.", err)
	//}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	//for {

	//pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	//_, err = clientset.CoreV1().ConfigMaps("default").Create(context.TODO(), &configMap, metav1.CreateOptions{})
	if _, err = clientset.CoreV1().ConfigMaps("default").Get(context.TODO(), "guardian", metav1.GetOptions{}); errors.IsNotFound(err) {
		_, err = clientset.CoreV1().ConfigMaps("default").Create(context.TODO(), &configMap, metav1.CreateOptions{})
	} else {
		_, err = clientset.CoreV1().ConfigMaps("default").Update(context.TODO(), &configMap, metav1.UpdateOptions{})
	}

	//cm, err := clientset.CoreV1().ConfigMaps("default").Create(context.TODO(), &cm, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("ConfigMap Error: %v\n", err)
		panic(err.Error())
	}

}
*/
func (l *learner) fetchConfig(w http.ResponseWriter, req *http.Request) {
	ns, sid, cmflag := processQuery(req.URL.Query())
	if sid == "" || ns == "" {
		http.Error(w, "Missing data", http.StatusBadRequest)
		return
	}

	fmt.Printf("Servicing fetchConfig of service id %s ns %s (cmflag %t)\n", sid, ns, cmflag)
	record := l.loadSession(ns, sid, cmflag)
	buf, err := json.Marshal(record.wsgate)
	if err != nil {
		fmt.Printf("Servicing fetchConfig error while Marshal %v\n", err)
	}
	w.Write(buf)
	/*
		data := new(spec.WsGate)
		data.Control.Consult = true
		data.Control.RequestsPerMinuete = 60
		data.Control.Block = false
		data.Configured = new(spec.Criteria)
		data.Configured.Req.AddTypicalVal()
		//data.Req.Url.Val.AddValExample("/")
		//fmt.Println("Url ", data.Req.Url.Val.Describe())
		//data.Req.Qs.Kv.WhitelistKnownKeys(map[string]string{"a": "4", "b": ""})
		//data.Req.Qs.Kv.WhitelistKnownKeys(map[string]string{"b": "abcdefg123 *?"})
		//data.Req.Qs.Kv.WhitelistByExample("aa", "+972-54-5445-321")
		//data.Req.Qs.Kv.WhitelistByExample("aaa", "123456")
		//data.Req.Qs.Kv.SetMandatoryKeys([]string{"a"})
		//fmt.Println("Qs ", data.Req.Qs.Kv.Describe())
		buf, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("Servicing fetchConfig error while Marshal %v\n", err)
		}
		w.Write(buf)
	*/
}

func main() {
	l := new(learner)
	l.services = make(map[string]*serviceRecord)
	l.namespaces = make(map[string]bool)
	l.kmgr.InitConfigs()
	//getGuardianClient()
	//setConfigMap()
	//data := new(spec.GuardianSpec)
	//data.Control = new(spec.Ctrl)
	//data.Control.Consult = true
	//data.Control.RequestsPerMinuete = 60
	//data.Control.Block = false
	//data.Configured = new(spec.Criteria)
	//data.Configured.Req.AddTypicalVal()
	//l.kmgr.SetCrd("default", "mytestservice", data)
	fmt.Printf("Starting guard-learner on port 80\n")
	http.HandleFunc("/config", l.fetchConfig)
	http.HandleFunc("/req", l.consultOnReq)
	http.HandleFunc("/pile", l.processPile)
	err := http.ListenAndServe(":80", nil)
	fmt.Printf("Failed to start %v\n", err)
}
