#!/bin/bash
#
# Spins up 4 Marathon apps and creates 4 App monitors
# To stop the monitors and remove the marathon apps: sample_apps.sh stop
#
#set -x
if [ -f /etc/redhat-release ]; then
    if ! rpm -qa | grep jq; then
        echo "Please install jq from epel"
        exit 1
    fi
fi

AS_URL="http://autoscaler.marathon.l4lb.thisdcos.directory"
MASTER_URL="http://leader.mesos"
NAMESPACE="myapp"
SU_USR=""
SU_PWD=""
TKN=$(curl -H "Accept: application/json" -H "Content-type: application/json" -X POST -d '{"uid":"'${SU_USR}'","password":"'${SU_PWD}'"}' http://master.mesos/acs/api/v1/auth/login | jq ".token" | xargs)

function start(){
    for i in {1..4}; do
        #dcos marathon app remove /${NAMESPACE}/test${i} &> /dev/null
        curl -H "Content-type: application/json" -H "Authorization: token=${TKN}" -X DELETE ${MASTER_URL}/service/marathon/v2/apps/${NAMESPACE}/test${i}
        sed "s/NAMESPACE/\/${NAMESPACE}/g" test/test${i}.json > /tmp/test${i}.json
        #dcos marathon app add /tmp/test${i}.json
        curl -H "Content-type: application/json" -H "Authorization: token=${TKN}" -X POST ${MASTER_URL}/service/marathon/v2/apps -d @/tmp/test${i}.json
        sleep 10s
    done

    echo "Creating sample application monitors"
    #cpu
    curl -X POST -d '{ "app_id": "/'${NAMESPACE}'/test1","max_cpu": 60,"min_cpu": 20,"max_mem": 90,"min_mem": 5,"method": "cpu","scale_factor": 1,"max_instances": 6,"min_instances": 2,"warm_up": 3,"cool_down": 3,"interval": 181}' ${AS_URL}/apps
    sleep 5s
    #mem
    curl -X POST -d '{"app_id": "/'${NAMESPACE}'/test2","max_cpu": 60,"min_cpu": 20,"max_mem": 75,"min_mem": 15,"method": "mem","scale_factor": 1,"max_instances": 5,"min_instances": 2,"warm_up": 3,"cool_down": 4,"interval": 183}' ${AS_URL}/apps
    sleep 5s
    #and
    curl -X POST -d '{"app_id": "/'${NAMESPACE}'/test3","max_cpu": 55,"min_cpu": 15,"max_mem": 80,"min_mem": 10,"method": "and","scale_factor": 1,"max_instances": 6,"min_instances": 2,"warm_up": 3,"cool_down": 3,"interval": 185}' ${AS_URL}/apps
    sleep 5s
    #or
    curl -X POST -d '{"app_id": "/'${NAMESPACE}'/test4","max_cpu": 70,"min_cpu": 50,"max_mem": 75,"min_mem": 50,"method": "or","scale_factor": 1,"max_instances": 5,"min_instances": 2,"warm_up": 3,"cool_down": 3,"interval": 187}' ${AS_URL}/apps

}

function stop(){
    for i in {1..4}; do
        curl -X DELETE -d '{"app_id": "/'${NAMESPACE}'/test'${i}'"}' ${AS_URL}/apps &> /dev/null
        #dcos marathon app remove /${NAMESPACE}/test${i} &> /dev/null
        curl -H "Content-type: application/json" -H "Authorization: token=${TKN}" -X DELETE ${MASTER_URL}/service/marathon/v2/apps/${NAMESPACE}/test${i}
        sleep 10s
    done
}

# Main


if  [ "${1}" == "stop" ]; then
    stop
else
    if ! curl ${AS_URL}; then
        echo "Please start autoscaler first."
        exit 1
    fi
    start
fi

echo "Done"
