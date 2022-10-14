# Guard
Guard is a security layer to protect your code engine Services.  Guard is developed as open source. The Guard project is presently being adopted by the [Knative Community](http://knative.dev/security-guard). Code Engine users are now able to try out Guard as part of their Code Engine services. Guard core component is the 

## Security 
[`guard-gate`](https:pkg/guard-gate) uses Micro-rules to monitor and potentialy block requests and/or responses to services. 
Micro-rules offer a fine grain filtering performed against each value delivered to/from the service.
By using Micro-rules, [`guard-gate`](https:pkg/guard-gate) makes it hard to deliver an exploit to be used against a vulnerability embedded as part of the service or its dependencies.

Using [`guard-gate`](https:pkg/guard-gate) will typically require an offender to build a dedicated delivery mechanism in order to explore options for detecting and exploiting service vulnerabilities. This leangthy process may need to be repeated for each service, as each service maintains a different set of micro-rules. As a result, an offender will not be able to use common statistical attack patterns.  

The user of [`guard-gate`](https:pkg/guard-gate) gains Situational Awareness both thanks to alerts about request/responses out of pattern and by the identification of indicators that the service is misused. Such indicators include longer than usual  service times and the list of external IP addresses appraoched by the service. 

Beyond, Situational Awareness, [`guard-gate`](https:pkg/guard-gate) enables blocking of out-of-pattern behaviours and the ability to react to potential attacks and/or to on-going attacks by introducing a fine-tune configrable security gate in front of the service.   

Overall the solution offers both visibility into the security of the service as well as the ability to monitor/block both known patterns and unknown patterns (using a zero day exploits).

## Solution components
This project adds:
1. A workload security gate named [`guard-gate`](https:pkg/guard-gate) implemented as a go package
1. A set of micro-rules used by [`guard-gate`](https:pkg/guard-gate) named *Guardian* implemenyedas a CRD or configmap
1. A learner service to auto learn the micro-rules in *Guardian* named [`guard-learner`](https:cmd/guard-learner)
1. A user interface web app to simplify manual configuration of micro-rules named [`guard-ui`](https:cmd/guard-ui) 

In addition, the project adds:
1. A go package to enable using [`guard-gate`](https:pkg/guard-gate) woth Knative Queue Proxy Option named [`qpoption`](https:pkg/qpoption)
2. A go package to enable testing Knative Queue Proxy Options named [`test-gate`](https:pkg/test-gate)

## Guard Gate
[`guard-gate`](https:pkg/guard-gate) can be loaded as a knative queue proxy option using [`qpoption`](https:pkg/qpoption)

Once loaded, it monitors the proxied requests and responses. 
If the proxy runs as a sidecar conatiner to the service, the pod network may also be monitored.

Note that [`guard-gate`](https:pkg/guard-gate) can also be used for more generic kubernetes use cases by loading:
- As a standalone reverse proxy, see for example: [`guard-rproxy`](https://github.com/IBM/workload-security-guard/tree/main/cmd/guard-rproxy)
- As an extension to any go proxy, for example by using: [`rtplugs`](https://github.com/IBM/go-security-plugs/tree/main/rtplugs).


## Guardian
[`guard-gate`](https:pkg/guard-gate) uses *Guardian* - a set of micro-rules that define the expected behaviour of the service.

*Guardian* may reside in a CRD (guardians.wsecurity.ibmresearch.com) under the name <servicename>.<namespace> or in a configmap under the name 'guardian-<servicename>'. If a *Guardian* is not found, [`guard-gate`](https:pkg/guard-gate) will look for a namespace-default *Guardian* as a starting point under the name  <ns>-<namespace> or in a configmap under the name 'guardian-<ns>-<namespace>'.  If a namespace-default *Guardian* is not found, [`guard-gate`](https:pkg/guard-gate) will use an empty set of micro-rules as a starting point and will set itself to work in auto-learning mode (See [`guard-gate`](https:pkg/guard-gate) for more details on the different working modes).

## Guard Learner 
[`guard-learner`](https:cmd/guard-learner) is a standalone service used to learn *Guardian* micro-srules based on inputs from instances of [`guard-gate`](https:pkg/guard-gate). [`guard-learner`](https:cmd/guard-learner) stores the *Guardian* as a CRD (guardians.wsecurity.ibmresearch.com) under the name <servicename>.<namespace> or in a configmap under the name 'guardian-<servicename>'.

## Guard User Interface 
Although *Guardian* CRDs and Configmaps can be controled directly via kubectl. An optional [`guard-ui`](https:cmd/guard-ui) is offered to simplify and clarify the micro-rules. 

