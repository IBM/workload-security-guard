package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	//"unicode"

	"github.com/IBM/go-security-plugs/rtplugs"
	_ "github.com/IBM/workload-security-guard/pkg/guardgate"
	"go.uber.org/zap"
)

//func test(r rune) {
//	for name, table := range unicode.Scripts {
//		if unicode.Is(table, r) {
//			fmt.Printf("Found %s\n", name)
//		}
//	}
//}

// Eample of a Reverse Proxy using plugs
func main() {
	//loadUnicodeTable()

	//str := "אבגדהוזחטיאבגדהוזחטיאבגדהוזחטיאבגדהוזחטיאבגדהוזחטיאבגדהוזחטיאבגדהוזחטיאבגדהוזחטיאבגדהוזחטיאבגדהוזחטי"
	//runes := []rune(str)
	//start := time.Now()
	//for i := 0; i < len(runes); i++ {
	//		test(runes[i])
	//	}

	//	test([]rune("ג")[0])
	//	elapsed := time.Since(start)
	//	fmt.Printf("took %s", elapsed)
	// Sleep some, such theb  gaurd can start before you
	time.Sleep(3 * time.Second)
	logger, _ := zap.NewDevelopment()
	//logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	url, err := url.Parse("http://127.0.0.1:8082")
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	var h http.Handler = proxy

	// Hook using RoundTripper
	os.Setenv("WSGATE_GUARD_URL", "http://127.0.0.1:8888")
	os.Setenv("RTPLUGS", "wsgate")
	os.Setenv("SERVING_NAMESPACE", "default")
	os.Setenv("SERVING_SERVICE", "testservice")
	os.Setenv("CMNAME", "true")
	rt := rtplugs.New(log)
	if rt != nil {
		defer rt.Close()
		proxy.Transport = rt.Transport(proxy.Transport)
	}

	http.Handle("/", h)
	log.Fatal(http.ListenAndServe("localhost:8081", nil))
}
