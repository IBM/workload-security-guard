package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"path/filepath"

	spec "github.com/IBM/workload-security-guard/pkg/apis/wsecurity/v1"
	guardianclientset "github.com/IBM/workload-security-guard/pkg/generated/clientset/guardians"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func processPile(w http.ResponseWriter, req *http.Request) {
	// Add security to ensure that only gate can use this interface!
	// Check that the source is from the local 10.*.*.* range
	// Check that the
	var pile spec.ReqPile
	query := req.URL.Query()
	sidSlice := query["sid"]
	nsSlice := query["ns"]
	fmt.Printf("Servicing processPile %v\n", query)
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
	data.Contigured.Req.AddTypicalVal()

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

func setCrd(sid string, wsGate *spec.WsGate) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	guardianClient, err := guardianclientset.NewForConfig(cfg)
	if err != nil {
		panic(err.Error())
	}

	gClient := guardianClient.WsecurityV1()

	var g *spec.Guardian

	g, err = gClient.Guardians("default").Get(context.TODO(), sid, metav1.GetOptions{})
	if err == nil {
		fmt.Printf("guardian %v\n", g)
		g.Name = sid
		g.Spec = (*spec.GuardianSpec)(wsGate)
		g, err = gClient.Guardians("default").Update(context.TODO(), g, metav1.UpdateOptions{})
		if err != nil {
			fmt.Printf("update err %v\n", err)
		}
	} else {
		fmt.Printf("get err %v\n", err)
		g.Name = sid
		g.Spec = (*spec.GuardianSpec)(wsGate)
		g, err = gClient.Guardians("default").Create(context.TODO(), g, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("update err %v\n", err)
		}
	}
}

func setConfigMap() {
	data := new(spec.WsGate)
	data.Control.Consult = true
	data.Control.RequestsPerMinuete = 60
	//data.ForceAllow = true
	data.Contigured.Req.AddTypicalVal()

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
	fmt.Printf("My Marshal\n %s\n", data.Contigured.Req.Marshal(4))

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
	data.Contigured.Req.AddTypicalVal()
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
	//setConfigMap()
	data := new(spec.WsGate)
	data.Control.Consult = true
	data.Control.RequestsPerMinuete = 60
	data.Control.Block = false
	data.Contigured.Req.AddTypicalVal()
	setCrd("myservice.mynamepsace", data)
	fmt.Printf("Starting server on port 8888\n")
	http.HandleFunc("/config", fetchConfig)
	http.HandleFunc("/req", consultOnReq)
	http.HandleFunc("/pile", processPile)
	err := http.ListenAndServe("localhost:8888", nil)
	fmt.Printf("Failed to start %v\n", err)

}
