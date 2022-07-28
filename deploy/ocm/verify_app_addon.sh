#!/bin/bash

set -o nounset
set -o pipefail

waitForRes() {
    FOUND=1
    MINUTE=0
    resKinds=$1
    resName=$2
    resNamespace=$3
    ignore=$4
    running="\([0-9]\+\)\/\1"
    printf "\n#####\nWait for ${resNamespace}/${resName} to reach running state (4min).\n"
    while [ ${FOUND} -eq 1 ]; do
        # Wait up to 3min, should only take about 20-30s
        if [ $MINUTE -gt 180 ]; then
            echo "Timeout waiting for the ${resNamespace}\/${resName}."
            echo "List of current resources:"
            oc get ${resKinds} -n ${resNamespace} ${resName}
            echo "You should see ${resNamespace}/${resName} ${resKinds}"
            if [ "${resKinds}" == "pods" ]; then
                oc describe deployments -n ${resNamespace} ${resName}
            fi
            exit 1
        fi
        if [ "$ignore" == "" ]; then
            operatorRes=`oc get ${resKinds} -n ${resNamespace} | grep ${resName}`
        else
            operatorRes=`oc get ${resKinds} -n ${resNamespace} | grep ${resName} | grep -v ${ignore}`
        fi
        if [[ $(echo $operatorRes | grep "${running}") ]]; then
            echo "* ${resName} is running"
            break
        elif [ "$operatorRes" == "" ]; then
            operatorRes="Waiting"
        fi
        echo "* STATUS: $operatorRes"
        sleep 5
        (( MINUTE = MINUTE + 5 ))
    done
}

checkLeaseOutput() {
    MINUTE=0
    while [ true ]; do
        # Wait up to 3min
        if [ $MINUTE -gt 180 ]; then
            echo "3 minutes reached."
            exit 0
        fi
        $KUBECTL get lease -n open-cluster-management-agent-addon application-manager   -o yaml

        sleep 60
        (( MINUTE = MINUTE + 60 ))
    done
}

KUBECTL=${KUBECTL:-kubectl}


echo "############  access cluster1"
$KUBECTL config use-context kind-cluster1

waitForRes "pods" "application-manager" "open-cluster-management-agent-addon" ""

APP_ADDON_POD=$($KUBECTL get pods -n open-cluster-management-agent-addon |grep application-manager |awk -F' ' '{print $1}')

sleep 60

# output the application manager pod log
$KUBECTL logs -n open-cluster-management-agent-addon $APP_ADDON_POD

# output lease result twice, make sure the lease is updated every 1 min
checkLeaseOutput
