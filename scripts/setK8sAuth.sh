#!/usr/bin/env bash 
SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
projfullpath="$(cd ${SCRIPT_ROOT}; pwd)"
#echo "SCRIPT_ROOT $SCRIPT_ROOT"
#GOBIN="$(go env GOBIN)"
#gobin="${GOBIN:-$(go env GOPATH)/bin}"
#gosrc="${GOBIN:-$(go env GOPATH)/src}"
#boilerplate="${projfullpath}/scripts/boilerplate.go.txt"
#proj="github.com/IBM/workload-security-guard"

kubectl apply -f "${projfullpath}/deploy/role.yaml"
kubectl apply -f "${projfullpath}/deploy/serviceAccount.yaml"
kubectl apply -f "${projfullpath}/deploy/roleBinding.yaml"

