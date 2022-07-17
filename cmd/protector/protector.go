package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/IBM/go-security-plugs/rtplugs"
	_ "github.com/IBM/workload-security-guard/pkg/wsgate"
	"go.uber.org/zap"
)

type BaseRoundTrip struct {
	next http.RoundTripper // the next roundtripper
}

func (rt *BaseRoundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Printf("BaseRoundTrip req.Host %s\n", req.Host)
	fmt.Printf("BaseRoundTrip req.URL.Host %s\n", req.URL.Host)

	req.Host = "" // req.URL.Host
	return rt.next.RoundTrip(req)
}

func (rt *BaseRoundTrip) Transport(t http.RoundTripper) http.RoundTripper {
	if t == nil {
		t = http.DefaultTransport
	}
	rt.next = t
	return rt
}

// Eample of a Reverse Proxy using plugs
func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()
	log.Info("Protector Initializing")
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		log.Infof("ENV %s: %s\n", pair[0], pair[1])
	}
	log.Info("ENV________DONE")
	serviceTcp := os.Getenv("SERVICETCP")

	if serviceTcp == "" {
		serviceUrl := os.Getenv("SERVICEURL")
		if serviceUrl == "" {
			// default
			serviceUrl = "http://127.0.0.1:80"
		}
		log.Infof("protector serving serviceUrl: %s", serviceUrl)
		parsedUrl, err := url.Parse(serviceUrl)
		if err != nil {
			log.Fatalf("Filed to parse http service url: %s", err.Error())
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(parsedUrl)

		//var h http.Handler = proxy

		// Hook using BaseRoundTripper
		var brt BaseRoundTrip
		proxy.Transport = brt.Transport(proxy.Transport)

		// Hook using RoundTripper
		rt := rtplugs.New(log)
		if rt != nil {
			defer rt.Close()
			proxy.Transport = rt.Transport(proxy.Transport)
		}

		http.Handle("/", proxy)
		err = http.ListenAndServe(":22000", nil)
		log.Fatalf("Failed to open http local service: %s", err.Error())
	} else {
		listener, err := net.Listen("tcp", ":22000")
		if err != nil {
			log.Fatalf("Failed to open tcp local service: %s", err.Error())
			return
		}
		for {
			incoming, err := listener.Accept()
			log.Info("New connection", incoming.RemoteAddr())
			if err != nil {
				log.Info("error accepting connection", err)
				continue
			}
			go func() {
				defer incoming.Close()
				outgoing, err := net.Dial("tcp", serviceTcp)
				if err != nil {
					log.Info("error dialing remote addr", err)
					return
				}
				defer outgoing.Close()
				closer := make(chan struct{}, 2)
				go copy(closer, outgoing, incoming)
				go copy(closer, incoming, outgoing)
				<-closer
				log.Info("Connection complete", incoming.RemoteAddr())
			}()
		}
	}
}

func copy(closer chan struct{}, dst io.Writer, src io.Reader) {
	_, _ = io.Copy(dst, src)
	closer <- struct{}{} // connection is closed, send signal to stop proxy
}
