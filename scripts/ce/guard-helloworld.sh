#!/usr/bin/env bash

export SERVICE_NAME="myapp"

if [ -z "$SERVICE_NAME" ]
then
      echo "service-name must be set"
      exit
fi

export NAMESPACE=`ibmcloud ce project current -o json|jq -r .kube_config_context`
if [ -z "$NAMESPACE" ]
then
      echo "cant find current namespace"
      exit
fi


# get the guard-learner url
export GUARD_URL=`ibmcloud ce application get -n guard-learner -o url`

if [ -z "$GUARD_URL" ]
then
      echo "guard-learner is required"
      exit
fi


# create the applicatin that need to be protected - use any image and its respective exposed port here
# -min 1 is advised but not mandatory
echo "Creating a new code engine service named '${SERVICE_NAME}' (no access outside the project)"
ibmcloud ce application create -n ${SERVICE_NAME} -v project -min 1 -p 8080 -i icr.io/codeengine/hello
export SERVICE_URL=`ibmcloud ce application get -n ${SERVICE_NAME} -o url`

if [ -z "$SERVICE_URL" ]
then
      echo "guard-learner is required"
      exit
fi

# create the guard to protect the application
# -min 1 is advised but not mandatory
echo ""
echo "Creating a security guard to protect the new code engine service (public access)"
echo "The protected service name '${SERVICE_NAME}' namespace '${NAMESPACE}' url '${SERVICE_URL}'"
echo "The guard learner url is '${GUARD_URL}'"
ibmcloud ce application create -n ${SERVICE_NAME}-guard --min 1 -p 22000 \
        -e GUARD_URL=${GUARD_URL} \
        -e SERVICE_URL=${SERVICE_URL} \
        -e SERVICE_NAME=${SERVICE_NAME} \
        -e NAMESPACE=${NAMESPACE} \
        -e USE_CONFIGMAP=true \
        -i ghcr.io/ibm/workload-security-guard/guard-rproxy