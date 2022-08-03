
#Guardian role
#kubectl apply -f kube/role.yaml
#kubectl apply -f kube/serviceAccount.yaml
#kubectl apply -f kube/roleBinding.yaml
#kubectl apply -f namespace.yml
#kubectl apply -f service.yml

#kubectl get secret all-icr-io -n default -o yaml | sed 's/default/knative-guard/g' | kubectl create -n knative-guard -f -

export KO_DOCKER_REPO=us.icr.io/knat
ko apply -Rf deployment.yml

#bx login --apikey @~/keys/apikey.json
#ibmcloud cr login
#docker build  --tag guardian . 
#docker tag guardian us.icr.io/dev_sec_ops/guardian:latest
#docker push us.icr.io/dev_sec_ops/guardian:latest
#kubectl rollout restart deployment guardian -n knative-guardian

#CRDs
#kubectl apply -f kube/Gates.yaml
#kubectl apply -f kube/Guardians.yaml
#kubectl apply -f kube/envoy.gate.yaml

# Support for reading images from IBM Container Registry
#kubectl get secret all-icr-io -n default -o yaml | sed 's/default/tests/g' | kubectl create -n tests -f -
#kubectl get secret all-icr-io -n default -o yaml | sed 's/default/knative-guard/g' | kubectl create -n knative-guard -f -
#kubectl get secret all-icr-io -n default -o yaml | sed 's/default/knative-serving-ingress/g' | kubectl create -n knative-serving-ingress -f -
#kubectl patch ServiceAccount default -n tests -p '{"imagePullSecrets": [{"name": "all-icr-io"}]}'
#kubectl patch ServiceAccount default -n knative-guard -p '{"imagePullSecrets": [{"name": "all-icr-io"}]}'
#kubectl patch ServiceAccount default -n knative-serving-ingress -p '{"imagePullSecrets": [{"name": "all-icr-io"}]}'