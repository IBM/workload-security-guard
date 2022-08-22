#!/usr/bin/env bash
SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
projfullpath="$(cd ${SCRIPT_ROOT}; pwd)"

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
kubectl apply -f "${projfullpath}/deploy/role.yaml"
kubectl apply -f "${projfullpath}/deploy/serviceAccount.yaml"
kubectl apply -f "${projfullpath}/deploy/roleBinding.yaml"
#CRDs
#kubectl apply -f "${projfullpath}/deploy/Gates.yaml"
#kubectl apply -f "${projfullpath}/deploy/Guardians.yaml"
#kubectl apply -f "${projfullpath}/deploy/envoy.gate.yaml"

# Support for reading images from IBM Container Registry
kubectl get secret all-icr-io -n default -o yaml | sed 's/default/knative-serving/g' | kubectl create -n knative-serving -f -
kubectl patch ServiceAccount default -n knative-serving -p '{"imagePullSecrets": [{"name": "all-icr-io"}]}'
kubectl patch ServiceAccount controller -n knative-serving -p '{"imagePullSecrets": [{"name": "all-icr-io"}]}'
