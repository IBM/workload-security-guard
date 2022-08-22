package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"runtime/debug"

	pi "github.com/IBM/go-security-plugs/pluginterfaces"
	_ "github.com/IBM/workload-security-guard/pkg/test-gate"
	"knative.dev/pkg/signals"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type config struct {
	Sender      string `split_words:"true" required:"false"`
	Response    string `split_words:"true" required:"false"`
	ServiceName string `split_words:"true" required:"true"`
	Namespace   string `split_words:"true" required:"true"`
	ServiceUrl  string `split_words:"true" required:"true"`
	LogLevel    string `split_words:"true" required:"false"`
}

type GuardGate struct {
	nextRoundTripper http.RoundTripper // the next roundtripper
	securityPlug     pi.RoundTripPlug
}

func (p *GuardGate) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			pi.Log.Warnf("Recovered from panic during RoundTrip! Recover: %v\n", recovered)
			pi.Log.Debugf("Stacktrace from panic: \n %s\n" + string(debug.Stack()))
			err = errors.New("paniced during RoundTrip")
			resp = nil
		}
	}()
	req.Host = "" // req.URL.Host

	if req, err = p.securityPlug.ApproveRequest(req); err == nil {
		if resp, err = p.nextRoundTripper.RoundTrip(req); err == nil {
			resp, err = p.securityPlug.ApproveResponse(req, resp)
		}
	}
	if err != nil {
		pi.Log.Debugf("%s: returning error %v", p.securityPlug.PlugName(), err)
		resp = nil
	}
	return
}

func (p *GuardGate) Transport(t http.RoundTripper) http.RoundTripper {
	if t == nil {
		t = http.DefaultTransport
	}
	p.nextRoundTripper = t
	return p
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel

	}
}

func createLogger(logLevel string) *zap.SugaredLogger {
	rawJSON := []byte(`{
		"level": "info",
		"encoding": "json",
		"outputPaths": ["stdout"],
		"development": false,
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase"
		}
	  }`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	cfg.Level = zap.NewAtomicLevelAt(getLogLevel(logLevel))
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger.Sugar()
}

// Eample of a Reverse Proxy using plugs
func main() {
	var env config
	if err := envconfig.Process("", &env); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to process environment: %s\n", err.Error())
		os.Exit(1)
	}

	plugConfig := make(map[string]string)
	plugConfig["response"] = env.Response
	plugConfig["sender"] = env.Sender

	log := createLogger(env.LogLevel)
	defer log.Sync()
	pi.Log = log

	log.Infof("guard-proxy serving serviceName: %s, namespace: %s, serviceUrl: %s", env.ServiceName, env.Namespace, env.ServiceUrl)
	parsedUrl, err := url.Parse(env.ServiceUrl)
	if err != nil {
		log.Fatalf("Failed to parse serviceUrl: %s", err.Error())
	}
	proxy := httputil.NewSingleHostReverseProxy(parsedUrl)

	// Hook using RoundTripper

	securityPlug := pi.RoundTripPlugs[0]
	securityPlug.Init(signals.NewContext(), plugConfig, env.ServiceName, env.Namespace, log)
	defer securityPlug.Shutdown()

	var gateGaurd GuardGate
	gateGaurd.securityPlug = securityPlug
	proxy.Transport = gateGaurd.Transport(proxy.Transport)

	http.Handle("/", proxy)
	log.Infof("Creating Reverse Proxy on port 8081")
	err = http.ListenAndServe(":8081", nil)
	log.Fatalf("Failed to open http local service: %s", err.Error())
}