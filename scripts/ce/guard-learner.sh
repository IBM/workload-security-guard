#!/usr/bin/env bash

# create the guard-learner - need only one such instance per project
ibmcloud ce application create -n guard-learner -v project --min 1 --max 1 -p 8888 -i ghcr.io/ibm/workload-security-guard/guard-learner
