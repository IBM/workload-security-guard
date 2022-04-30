package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"

	spec "github.com/IBM/workload-security-guard/pkg/apis/wsecurity/v1"
	guardianclientset "github.com/IBM/workload-security-guard/pkg/generated/clientset/guardians"
	v1 "github.com/IBM/workload-security-guard/pkg/generated/clientset/guardians/typed/wsecurity/v1"

	pi "github.com/IBM/go-security-plugs/pluginterfaces"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type guardianui struct {
	//namespace string
	//serviceId string
	gClient v1.WsecurityV1Interface
	//wsGate    *spec.WsGate
}

var gui *guardianui

func (gui *guardianui) setCrd(namespace string, serviceId string, guardianSpec *spec.GuardianSpec) bool {
	var g *spec.Guardian
	var err error

	fmt.Printf("setCrd %v\n", guardianSpec.Control)

	fmt.Print((*spec.WsGate)(guardianSpec).Marshal(0))
	g, err = gui.gClient.Guardians(namespace).Get(context.TODO(), serviceId, metav1.GetOptions{})
	if err == nil {
		fmt.Printf("setCrd: guardian read succesful %v\n", g.Spec.Control)
		g.Name = serviceId
		g.Spec = guardianSpec
		_, err = gui.gClient.Guardians(namespace).Update(context.TODO(), g, metav1.UpdateOptions{})
		if err != nil {
			fmt.Printf("setCrd: update err %v\n", err)
			return false
		}
		fmt.Printf("setCrd: guardian update succesfull %v\n", g.Spec.Control)
	} else {
		fmt.Printf("setCrd: guardian read err %v\n", err)
		g.Name = serviceId
		g.Spec = guardianSpec
		_, err = gui.gClient.Guardians(namespace).Create(context.TODO(), g, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("setCrd: guardian create err %v\n", err)
			return false
		}
		fmt.Printf("setCrd: guardian create succesfull %v\n", g)
	}
	return true
}

func (gui *guardianui) readCrd(namespace string, serviceId string) *spec.GuardianSpec {
	g, err := gui.gClient.Guardians(namespace).Get(context.TODO(), serviceId, metav1.GetOptions{})
	if err != nil {
		pi.Log.Infof("Err during get %s.%s: %s", serviceId, namespace, err.Error())
		//panic(fmt.Sprintf("No Guardian! for %s.%s", serviceId, namespace))
		return nil
	}
	pi.Log.Infof("Found guardian %s.%s", serviceId, namespace)
	fmt.Printf("Guardian received %v\n", g.Spec)

	fmt.Print((*spec.WsGate)(g.Spec).Marshal(0))
	return g.Spec
}

func (gui *guardianui) initCrd() {
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

	gui.gClient = guardianClient.WsecurityV1()
}

func setGuadian(w http.ResponseWriter, r *http.Request) {
	namespace := mux.Vars(r)["namespace"]
	service := mux.Vars(r)["service"]
	fmt.Printf("Guardian Set %s in namespace: %s\n", service, namespace)

	b, err1 := io.ReadAll(r.Body)
	if err1 != nil {
		fmt.Printf("Failed to read Body\n")
		return
	}
	str := string(b)
	fmt.Printf("Received: %s\n", str)
	var g spec.GuardianSpec

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.Unmarshal([]byte(str), &g)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	/*
		if g.Configured == nil {
			g.Configured = new(spec.Critiria)
		}
		g.Configured.Normalize()
	*/
	//Serr := json.NewDecoder(r.Body).Decode(&g)
	success := gui.setCrd(namespace, service, &g)
	// Do something with the Person struct...
	fmt.Printf("setGuadian: %v\n", success)
	resp := make(map[string]string)
	if success {
		w.WriteHeader(http.StatusOK)
		resp["message"] = "Resource was updated"
	} else {
		w.WriteHeader(http.StatusNotFound)
		resp["message"] = "Resource was not updated"
	}
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func getGuadian(w http.ResponseWriter, r *http.Request) {
	namespace := mux.Vars(r)["namespace"]
	service := mux.Vars(r)["service"]
	fmt.Printf("Guardian Get %s in namespace: %s\n", service, namespace)
	g := gui.readCrd(namespace, service)
	if g == nil {
		w.WriteHeader(http.StatusNotFound)
		resp := make(map[string]string)
		resp["message"] = "Resource Not Found"
		jsonResp, _ := json.Marshal(resp)
		w.Write(jsonResp)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

func main() {
	gui = new(guardianui)
	gui.initCrd()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/guardian/{namespace}/{service}", setGuadian).Methods("POST")
	router.HandleFunc("/guardian/{namespace}/{service}", getGuadian).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/build")))

	fmt.Println("Guardian App v0.01")
	log.Fatal(http.ListenAndServe(":9000", router))
}
