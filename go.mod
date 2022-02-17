module github.com/IBM/workload-security-guard

go 1.16

require (
	github.com/google/uuid v1.1.2
	go.uber.org/zap v1.19.1
	k8s.io/api v0.22.5
	k8s.io/apimachinery v0.22.5
	k8s.io/client-go v0.22.5
	k8s.io/klog/v2 v2.30.0
)

require github.com/IBM/go-security-plugs v1.1.1-0.20220217135704-2792fac15905
