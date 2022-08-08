#!/usr/bin/env bash

export SERVICE_NAME="myapp"


# clear previous runs
ibmcloud ce application delete -n ${SERVICE_NAME}
ibmcloud ce application delete -n ${SERVICE_NAME}-guard
ibmcloud ce application delete -n guard-learner
