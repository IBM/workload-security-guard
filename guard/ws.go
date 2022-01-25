package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IBM/workload-security-guard/spec"
)

/*
type wsgateConfig struct {
	QsKeys []spec.Minmax
	Tbd    []spec.Minmax
	Level  int
}*/

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
	data.Req.Url.Val.AddValExample("/")
	fmt.Println("Url ", data.Req.Url.Val.Describe())
	data.Req.Qs.Kv.WhitelistKnownKeys(map[string]string{"a": "4", "b": ""})
	data.Req.Qs.Kv.WhitelistKnownKeys(map[string]string{"b": "abcdefg123 *?"})
	data.Req.Qs.Kv.WhitelistByExample("aa", "+972-54-5445-321")
	data.Req.Qs.Kv.WhitelistByExample("aaa", "123456")
	data.Req.Qs.Kv.SetMandatoryKeys([]string{"a"})
	data.Req.Headers.AddTypicalVal()
	fmt.Println("Qs ", data.Req.Qs.Kv.Describe())
	buf, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Servicing fetchConfig error while Marshal %v\n", err)
	}
	//fmt.Printf("buf %s\n", string(buf))
	w.Write(buf)
}

func main() {
	fmt.Printf("Starting server on port 8888\n")
	http.HandleFunc("/fetchConfig", fetchConfig)
	err := http.ListenAndServe("localhost:8888", nil)
	fmt.Printf("Failed to start %v\n", err)

}
