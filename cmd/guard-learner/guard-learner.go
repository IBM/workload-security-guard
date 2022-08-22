package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	spec "github.com/IBM/workload-security-guard/pkg/apis/wsecurity/v1alpha1"
	guardkubemgr "github.com/IBM/workload-security-guard/pkg/guard-kubemgr"

	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

type serviceRecord struct {
	wsgate *spec.GuardianSpec
	pile   spec.SessionDataPile
}

type learner struct {
	services   map[string]*serviceRecord
	kmgr       guardkubemgr.KubeMgr
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
	fmt.Printf("processPile starts\n")
	var pile spec.SessionDataPile
	ns, sid, cmflag := processQuery(req.URL.Query())
	if sid == "" || ns == "" {
		fmt.Printf("processPile Missing data")
		http.Error(w, "processPile Missing data", http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(req.Body).Decode(&pile)
	if err != nil {
		fmt.Printf("processPile error\n%s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	{
		manifestJson, _ := json.MarshalIndent(pile, "", "  ")
		fmt.Println(string(manifestJson))
	}
	fmt.Printf("processPile about to load")
	record := l.loadSession(ns, sid, cmflag)

	record.pile.Merge(&pile)

	learn := new(spec.SessionDataConfig)

	learn.Learn(&record.pile)
	if record.wsgate.Learned != nil {
		learn.Fuse(record.wsgate.Learned)
	}
	record.wsgate.Learned = learn
	record.wsgate.Learned.Active = true
	fmt.Printf("About to setCrd record.wsgate.Control %v\n", record.wsgate.Control)

	if cmflag {
		ret := l.kmgr.SetCm(ns, sid, record.wsgate)
		fmt.Printf("Pile SetCm returned %s\n", ret)
	} else {
		ret := l.kmgr.SetCrd(ns, sid, record.wsgate)
		fmt.Printf("Pile setCrd returned %s\n", ret)
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
		fmt.Printf("Servicing fetchConfig error while JSON Marshal %v\n", err)
	}
	w.Write(buf)
}

func main() {
	l := new(learner)
	l.services = make(map[string]*serviceRecord)
	l.namespaces = make(map[string]bool)
	l.kmgr.InitConfigs()
	fmt.Printf("Starting guard-learner on port 8888\n")
	http.HandleFunc("/config", l.fetchConfig)
	http.HandleFunc("/req", l.consultOnReq)
	http.HandleFunc("/pile", l.processPile)
	err := http.ListenAndServe(":8888", nil)
	fmt.Printf("Failed to start %v\n", err)
}
