package v1

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReqProfile_Profile(t *testing.T) {
	var req *http.Request = httptest.NewRequest("GET", "/", nil)
	var cip net.IP = net.IPv4(1, 2, 3, 5)

	t.Run("ReqProfile", func(t *testing.T) {
		rp := new(ReqProfile)
		rc := new(ReqConfig)
		rPile := new(ReqPile)
		rp.Profile(req, cip)
		rp.Marshal(0)
		rPile.Clear()
		rPile.Add(rp)
		rPile.Marshal(0)
		rPile.Append(rPile)
		rc.AddTypicalVal()
		rc.Reconcile()
		rc.Marshal(0)
		rc.Decide(rp)
		rc.Learn(rPile)
		rc.Decide(rp)
		rc.Normalize()
		rc.Merge(rc)
		rPile.Clear()

	})

}

func TestRespProfile_Profile(t *testing.T) {
	var resp *http.Response = CreateRespone()

	t.Run("RespProfile", func(t *testing.T) {
		rp := new(RespProfile)
		rc := new(RespConfig)
		rPile := new(RespPile)
		rp.Profile(resp)
		rp.Marshal(0)
		rPile.Clear()
		rPile.Add(rp)
		rPile.Marshal(0)
		rPile.Append(rPile)
		rc.AddTypicalVal()
		rc.Reconcile()
		rc.Marshal(0)
		rc.Decide(rp)
		rc.Learn(rPile)
		rc.Normalize()
		rc.Merge(rc)
		rPile.Clear()
	})

}

func CreateRespone() *http.Response {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"value":"fixed"}`))
	}))
	defer server.Close()

	request, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	request.Header.Add("Accept", "application/json")
	client := &http.Client{}

	response, _ := client.Do(request)
	return response
}
