#bx login --apikey @~/keys/apikey.json
#ibmcloud cr login
#ibmcloud oc cluster config -c SROS2 
#KEY=`jq .apikey ~/keys/apikey.json`
#oc login -u davidh -p $KEY
#docker build  --tag guard-service
#docker tag guard-service us.icr.io/dev_sec_ops/guard-service:latest
#docker push us.icr.io/dev_sec_ops/guard-service:latest
#kubectl apply -f kube/Gates.yaml
#kubectl apply -f kube/envoy.gate.yaml
#kubectl get secret all-icr-io -n default -o yaml | sed 's/default/knative-serving-ingress/g' | kubectl create -n knative-serving-ingress -f -
#kubectl patch ServiceAccount default -n knative-serving-ingress -p '{"imagePullSecrets": [{"name": "all-icr-io"}]}'
#kubectl get secret all-icr-io -n default -o yaml | sed 's/default/tests/g' | kubectl create -n tests -f -
#kubectl patch ServiceAccount default -n tests -p '{"imagePullSecrets": [{"name": "all-icr-io"}]}'

kubectl apply -f kube/role.yaml
kubectl apply -f kube/serviceAccount.yaml
kubectl apply -f kube/roleBinding.yaml
kubectl apply -f namespace.yml
kubectl apply -f service.yml
kubectl apply -f kube/Guardians.yaml

# Support for reading images from IBM Container Registry
kubectl get secret all-icr-io -n default -o yaml | sed 's/default/knative-guard/g' | kubectl create -n knative-guard -f -
kubectl patch ServiceAccount default -n knative-guard -p '{"imagePullSecrets": [{"name": "all-icr-io"}]}'

export KO_DOCKER_REPO=us.icr.io/knat
ko apply -Rf deployment.yml

kubectl rollout restart deployment guard-service -n knative-guard
