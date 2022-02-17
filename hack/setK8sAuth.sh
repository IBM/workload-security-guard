#!/usr/bin/env bash 
SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
echo "SCRIPT_ROOT $SCRIPT_ROOT"
#GOBIN="$(go env GOBIN)"
#gobin="${GOBIN:-$(go env GOPATH)/bin}"
#gosrc="${GOBIN:-$(go env GOPATH)/src}"
projfullpath="$(cd ${SCRIPT_ROOT}; pwd)"
#boilerplate="${projfullpath}/hack/boilerplate.go.txt"
#proj="github.com/IBM/workload-security-guard"

kubectl apply -f "${projfullpath}/kube/role.yaml"
kubectl apply -f "${projfullpath}/kube/serviceAccount.yaml"
kubectl apply -f "${projfullpath}/kube/roleBinding.yaml"