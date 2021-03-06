bx login --apikey @~/keys/apikey.json
ibmcloud cr login
#ibmcloud oc cluster config -c SROS2 
#KEY=`jq .apikey ~/keys/apikey.json`
#oc login -u davidh -p $KEY
docker build  --tag guardian . 
docker tag guardian us.icr.io/dev_sec_ops/guardian:latest
docker push us.icr.io/dev_sec_ops/guardian:latest
kubectl rollout restart deployment guardian -n knative-guardian

#Guardian role
kubectl apply -f kube/role.yaml
kubectl apply -f kube/serviceAccount.yaml
kubectl apply -f kube/roleBinding.yaml
#CRDs
kubectl apply -f kube/Gates.yaml
kubectl apply -f kube/Guardians.yaml
kubectl apply -f kube/envoy.gate.yaml

# Support for reading images from IBM Container Registry
kubectl get secret all-icr-io -n default -o yaml | sed 's/default/tests/g' | kubectl create -n tests -f -
kubectl get secret all-icr-io -n default -o yaml | sed 's/default/knative-guardian/g' | kubectl create -n knative-guardian -f -
kubectl get secret all-icr-io -n default -o yaml | sed 's/default/knative-serving-ingress/g' | kubectl create -n knative-serving-ingress -f -
kubectl patch ServiceAccount default -n tests -p '{"imagePullSecrets": [{"name": "all-icr-io"}]}'
kubectl patch ServiceAccount default -n knative-guardian -p '{"imagePullSecrets": [{"name": "all-icr-io"}]}'
kubectl patch ServiceAccount default -n knative-serving-ingress -p '{"imagePullSecrets": [{"name": "all-icr-io"}]}'
