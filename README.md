# workload-security-guard
## Summary
This project adds a workload security gate to a kubernetes service. 
It is loaded as an extension to a go proxy using [`rtplugs`](https://github.com/IBM/go-security-plugs/rtplugs).
Once loaded, it monitors the proxied requests and responses. 
If the proxy runs a a sidecar conatiner, the pod can also be monitored.

A Guardian CRD incldues the Gate specifications and controlls the Gate behaviour.

[`guardui`](https:cmd/guardui) can be used to enable user interaction with the Guardian CRD


[`guard`](https:cmd/guard) can be used to auto learn the appropriate specifications ans drastically reduce or in some cases eliminate the required user interaction. 


## Security
The Gate makes it hard to deliver an exploit to be used against a vulnerability embedded as part of the service or its dependencies.
As a general rule, an attacker will be required to build a dedicated delivery mechanism to explore options for detecting and exploiting vulnerabilities for each service and will not be able to use common statistical attacks patterns.  

This is achieved thanks to a find grain filtering performed against each value delivered to the service.

Additional filtering enable identification of indicators that the service is misused.

Overall the solution offers both visibility into the security of the service as well as the ability to block both known patterns and unknown patterns (using a zero day exploit).

