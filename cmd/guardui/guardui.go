package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/gorilla/mux"

	spec "github.com/IBM/workload-security-guard/pkg/apis/wsecurity/v1"
	"github.com/IBM/workload-security-guard/pkg/wsgate"

	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

type guardianui struct {
	kmgr wsgate.Kubemgr
}

var gui *guardianui

func setGuadian(w http.ResponseWriter, r *http.Request) {
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
			g.Configured = new(spec.Critiria)
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

func getGuadian(w http.ResponseWriter, r *http.Request) {
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

func getCodeDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Dir(file)
}

func main() {
	gui = new(guardianui)
	//gui.initConfigs()
	gui.kmgr.InitConfigs()
	d := getCodeDir()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/guardian/{where}/{namespace}/{service}", setGuadian).Methods("POST")
	router.HandleFunc("/guardian/{where}/{namespace}/{service}", getGuadian).Methods("GET")
	path := filepath.Join(d, "frontend", "build")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir(path)))

	fmt.Println("Guardian App v0.01")
	log.Fatal(http.ListenAndServe(":9000", router))
}
