package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	spec "github.com/IBM/workload-security-guard/pkg/apis/wsecurity/v1"
	guardianclientset "github.com/IBM/workload-security-guard/pkg/generated/clientset/guardians"
	v1 "github.com/IBM/workload-security-guard/pkg/generated/clientset/guardians/typed/wsecurity/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type serviceRecord struct {
	wsgate *spec.WsGate
	pile   spec.Pile
}

var services = make(map[string]*serviceRecord)

func processPile(w http.ResponseWriter, req *http.Request) {
	// Add security to ensure that only gate can use this interface!
	// Check that the source is from the local 10.*.*.* range
	// Check that the
	var pile spec.Pile
	query := req.URL.Query()
	sidSlice := query["sid"]
	nsSlice := query["ns"]
	//fmt.Printf("Servicing processPile %v\n", query)
	if len(sidSlice) != 1 || len(nsSlice) != 1 {
		fmt.Printf("Servicing processPile missing data sid %d ns %d\n", len(sidSlice), len(nsSlice))
		return
	}
	sid := sidSlice[0]
	ns := nsSlice[0]
	if sid == "" || ns == "" {
		fmt.Printf("Servicing processPile missing data\n")
		http.Error(w, "Missing data", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(req.Body).Decode(&pile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("Servicing processPile of service id %s:%s %v\n", ns, sid, pile)
	record := loadSession(ns, sid)

	record.pile.Append(&pile)
	fmt.Printf("Record %s:%s %v\n", ns, sid, record.pile)
	fmt.Printf("Marshal Record\n______________\n%s\n", record.pile.Marshal())
	learn := new(spec.Critiria)
	learn.Learn(&record.pile)
	fmt.Printf("learn %v\n", learn)
	if len(record.wsgate.Learned) > 0 {
		learn.Merge(record.wsgate.Learned[0])
		record.wsgate.Learned[0] = learn
	} else {
		record.wsgate.Learned = append(record.wsgate.Learned, learn)
	}
	record.wsgate.Learned[0].Active = true
	fmt.Printf("record.wsgate.Learned %v\n", record.wsgate.Learned[0])
	fmt.Printf("Marshal Learned\n______________\n%s\n", record.wsgate.Learned[0].Marshal(0))
	updateCrd(ns, sid, record.wsgate)

	data := ""
	w.Write([]byte(data))
}

func consultOnReq(w http.ResponseWriter, req *http.Request) {
	// Add security to ensure that only gate can use this interface!
	// Check that the source is from the local 10.*.*.* range
	// Check that the
	query := req.URL.Query()
	sidSlice := query["sid"]
	nsSlice := query["ns"]
	fmt.Printf("Servicing consultOnReq %v\n", query)
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
	fmt.Printf("Servicing consultOnReq of service id %s:%s\n", ns, sid)
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
	data.Configured = new(spec.Critiria)
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
var gClient v1.WsecurityV1Interface

func getGuardianClient() {
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

func storeSession(ns string, sid string, wsgate *spec.WsGate) {
	if wsgate == nil {
		return
	}
	service := sid + "." + ns
	record, exists := services[service]
	if !exists {
		record = new(serviceRecord)
		record.pile.Clear()
		services[service] = record
	}
	record.wsgate = wsgate
}

func loadSession(ns string, sid string) *serviceRecord {
	service := sid + "." + ns
	record, exists := services[service]
	if exists {
		return record
	}

	// not cached
	wsgate := getCrd(ns, sid)
	if wsgate == nil {
		wsgate = getCrd("knative-serving", "guardian")
	}
	if wsgate == nil {
		fmt.Println("Guardian was not set!")
		wsgate = new(spec.WsGate)
	}

	storeSession(ns, sid, wsgate)
	return services[service]
}

func watchCrd() {
	for {
		watchCrdOnce()
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

func watchCrdOnce() {
	defer func() {
		if recovered := recover(); recovered != nil {
			fmt.Printf("Recovered from panic during watchCrdOnce! Recover: %v\n", recovered)
		}
	}()
	watcher, err := gClient.Guardians("").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("watchCrd err %v\n", err)
		return
	}
	ch := watcher.ResultChan()
	for {
		select {
		case event, ok := <-ch:
			if !ok {
				// the channel got closed, so we need to restart
				fmt.Printf("Kubernetes hung up on us, restarting event watcher")
				return
			}

			// handle the event
			fmt.Printf("event %v\n", event)

			switch event.Type {
			case watch.Added:
				g := event.Object.(*spec.Guardian)
				ns := g.ObjectMeta.Namespace
				sid := g.ObjectMeta.Name
				fmt.Printf("Service %s/%s added\n", ns, sid)
				storeSession(ns, sid, (*spec.WsGate)(g.Spec))
			case watch.Modified:
				g := event.Object.(*spec.Guardian)
				ns := g.ObjectMeta.Namespace
				sid := g.ObjectMeta.Name
				fmt.Printf("Service %s/%s modified\n", ns, sid)
				storeSession(ns, sid, (*spec.WsGate)(g.Spec))
			case watch.Deleted:
				g := event.Object.(*spec.Guardian)
				ns := g.ObjectMeta.Namespace
				sid := g.ObjectMeta.Name
				fmt.Printf("Service %s/%s deleted\n", ns, sid)
				delete(services, sid+"."+ns)
			case watch.Error:
				fmt.Printf("Error during watch %v\n", event.Object)
			}
		case <-time.After(10 * time.Minute):
			// deal with the issue where we get no events
			fmt.Printf("Timeout, restarting event watcher")
			return
		}
	}
}

func getCrd(ns string, sid string) *spec.WsGate {
	var g *spec.Guardian
	var err error

	fmt.Printf("getCrd guardian %s %s\n", ns, sid)

	g, err = gClient.Guardians(ns).Get(context.TODO(), sid, metav1.GetOptions{})
	if err == nil {
		fmt.Printf("getCrd err %v\n", err)
		return nil
	}
	return (*spec.WsGate)(g.Spec)
}

func updateCrd(ns string, sid string, wsGate *spec.WsGate) {
	var g *spec.Guardian
	var err error

	fmt.Printf("updateCrd guardian %s %s\n", ns, sid)
	g.Name = sid
	g.Spec = (*spec.GuardianSpec)(wsGate)
	_, err = gClient.Guardians(ns).Update(context.TODO(), g, metav1.UpdateOptions{})
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

/*
func setConfigMap() {
	data := new(spec.WsGate)
	data.Control.Consult = true
	data.Control.RequestsPerMinuete = 60
	//data.ForceAllow = true
	data.Configured = new(spec.Critiria)
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
func fetchConfig(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	sidSlice := query["sid"]
	fmt.Printf("Servicing fetchConfig %v\n", query)
	if len(sidSlice) != 1 {
		fmt.Printf("Servicing fetchConfig missing data %d\n", len(sidSlice))
		return
	}
	sid := sidSlice[0]
	if sid == "" {
		fmt.Printf("Servicing fetchConfig missing data\n")
		return
	}
	fmt.Printf("Servicing fetchConfig of service id %s\n", sid)
	data := new(spec.WsGate)
	data.Control.Consult = true
	data.Control.RequestsPerMinuete = 60
	data.Control.Block = false
	data.Configured = new(spec.Critiria)
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
}

func main() {
	getGuardianClient()
	watchCrd()
	//setConfigMap()
	data := new(spec.WsGate)
	data.Control.Consult = true
	data.Control.RequestsPerMinuete = 60
	data.Control.Block = false
	data.Configured = new(spec.Critiria)
	data.Configured.Req.AddTypicalVal()
	setCrd("default", "mytestservice", data)
	fmt.Printf("Starting server on port 8888\n")
	http.HandleFunc("/config", fetchConfig)
	http.HandleFunc("/req", consultOnReq)
	http.HandleFunc("/pile", processPile)
	err := http.ListenAndServe(":8888", nil)
	fmt.Printf("Failed to start %v\n", err)

}
