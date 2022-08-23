package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"

	spec "github.com/IBM/workload-security-guard/pkg/apis/wsecurity/v1"
	guardkubemgr "github.com/IBM/workload-security-guard/pkg/guard-kubemgr"

	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

type config struct {
	ServiceName     string `split_words:"true" required:"false"`
	Namespace       string `split_words:"true" required:"false"`
	UseConfigmap    bool   `split_words:"true" required:"false"`
	LockServiceName bool   `split_words:"true" required:"false"`
	LockNamespace   bool   `split_words:"true" required:"false"`
	LockConfigmap   bool   `split_words:"true" required:"false"`
	LogLevel        string `split_words:"true" required:"false"`
}

type guardianui struct {
	kmgr         guardkubemgr.Kubemgr
	currentSetup config
}

func (gui *guardianui) setGuadian(w http.ResponseWriter, r *http.Request) {
	where := mux.Vars(r)["where"]
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
			g.Configured = new(spec.Criteria)
		}
		g.Configured.Normalize()
	*/
	//Serr := json.NewDecoder(r.Body).Decode(&g)
	var ret string
	if where == "crd" {
		ret = gui.kmgr.SetCrd(namespace, service, &g)
	} else {
		ret = gui.kmgr.SetCm(namespace, service, &g)
	}
	// Do something with the Person struct...
	fmt.Printf("setGuadian: %v\n", ret)

	resp := make(map[string]string)
	if ret == "" {
		w.WriteHeader(http.StatusOK)
		resp["message"] = "Resource was updated"
	} else {
		w.WriteHeader(http.StatusNotFound)
		resp["message"] = ret
	}
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func (gui *guardianui) getGuadian(w http.ResponseWriter, r *http.Request) {
	where := mux.Vars(r)["where"]
	namespace := mux.Vars(r)["namespace"]
	service := mux.Vars(r)["service"]
	fmt.Printf("Guardian Get %s in namespace: %s\n", service, namespace)

	var g *spec.GuardianSpec

	if where == "crd" {
		g = gui.kmgr.ReadCrd(namespace, service)
	} else {
		g = gui.kmgr.ReadConfigMap(namespace, service)
	}
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

func (gui *guardianui) setup(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Guardian Setup\n")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gui.currentSetup)
}

func (gui *guardianui) initialize() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			fmt.Println("")
			fmt.Println("*-----------------------------------------------------------*")
			fmt.Println("* Unable to communicate with KubeAPI                        *")
			fmt.Println("*                                                           *")
			fmt.Println("* 1. Login to ibm cloud:                                    *")
			fmt.Println("*    > ibmcloud loin                                        *")
			fmt.Println("*                                                           *")
			fmt.Println("* 2. Connect to a code engine project:                      *")
			fmt.Println("*    > ibmcloud ce project select --name <ProjectName> -k   *")
			fmt.Println("*                                                           *")
			fmt.Println("* Then run:                                                 *")
			fmt.Println("*    > guard-ui                                             *")
			fmt.Println("*-----------------------------------------------------------*")
			fmt.Println("")
		}
	}()
	gui.kmgr.InitConfigs()
}

func getCodeDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Dir(file)
}

func main() {
	gui := new(guardianui)
	gui.initialize()

	if err := envconfig.Process("", &gui.currentSetup); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to process environment: %s\n", err.Error())
		os.Exit(1)
	}

	// path to index file when running from code
	d := getCodeDir()
	path := filepath.Join(d, "frontend/build")
	if _, err := os.Stat(path); err != nil {
		// path to index file when running from a container
		path = "/frontend"
	}
	fmt.Println("Guardian App v0.01")
	fmt.Println("Serving from", path)
	fmt.Printf("Setup Namespace %s ServiceName %s UseConfigmap %t LockConfigmap %t\n",
		gui.currentSetup.Namespace,
		gui.currentSetup.ServiceName,
		gui.currentSetup.UseConfigmap,
		gui.currentSetup.LockConfigmap)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/guardian/{where}/{namespace}/{service}", gui.setGuadian).Methods("POST")
	router.HandleFunc("/guardian/{where}/{namespace}/{service}", gui.getGuadian).Methods("GET")
	router.HandleFunc("/setup", gui.setup).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir(path)))

	fmt.Println("Services on port 9000")
	log.Fatal(http.ListenAndServe(":9000", router))
}
