package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/IBM/go-security-plugs/rtplugs"
	_ "github.com/IBM/workload-security-guard/wsgate"
	"go.uber.org/zap"
)

// Eample of a Reverse Proxy using plugs
func main() {
	// Sleep some, such theb  gaurd can start before you
	time.Sleep(3 * time.Second)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	url, err := url.Parse("http://127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	var h http.Handler = proxy

	// Hook using RoundTripper
	os.Setenv("WSGATE_GUARD_URL", "http://127.0.0.1:8888")
	os.Setenv("RTPLUGS", "wsgate")
	os.Setenv("SERVING_NAMESPACE", "mynamepsace")
	os.Setenv("SERVING_SERVICE", "myservice")
	rt := rtplugs.New(log)
	if rt != nil {
		defer rt.Close()
		proxy.Transport = rt.Transport(proxy.Transport)
	}

	http.Handle("/", h)
	log.Fatal(http.ListenAndServe("localhost:8081", nil))
}
