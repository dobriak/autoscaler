#!/bin/bash
#
# Spins up 4 Marathon apps and creates 4 App monitors
# To stop the monitors and remove the marathon apps: sample_apps.sh stop
#
set -x
AS_URL="http://localhost:8080"
function start(){
    for i in {1..4}; do
        dcos marathon app remove test${i}
        dcos marathon app add test/test${i}.json
        sleep 10s
    done

    echo "Creating sample application monitors"
    #cpu
    curl -X POST -d '{ "app_id": "/test1","max_cpu": 60,"min_cpu": 20,"max_mem": 90,"min_mem": 5,"method": "cpu","scale_factor": 1,"max_instances": 6,"min_instances": 2,"warm_up": 3,"cool_down": 3,"interval": 19}' ${AS_URL}/apps
    sleep 5s
    #mem
    curl -X POST -d '{"app_id": "/test2","max_cpu": 60,"min_cpu": 20,"max_mem": 75,"min_mem": 15,"method": "mem","scale_factor": 1,"max_instances": 5,"min_instances": 2,"warm_up": 3,"cool_down": 4,"interval": 21}' ${AS_URL}/apps
    sleep 5s
    #and
    curl -X POST -d '{"app_id": "/test3","max_cpu": 55,"min_cpu": 15,"max_mem": 80,"min_mem": 10,"method": "and","scale_factor": 1,"max_instances": 6,"min_instances": 2,"warm_up": 3,"cool_down": 3,"interval": 23}' ${AS_URL}/apps
    sleep 5s
    #or
    curl -X POST -d '{"app_id": "/test4","max_cpu": 70,"min_cpu": 50,"max_mem": 75,"min_mem": 50,"method": "or","scale_factor": 1,"max_instances": 5,"min_instances": 2,"warm_up": 3,"cool_down": 3,"interval": 25}' ${AS_URL}/apps

}

function stop(){
    for i in {1..4}; do
        curl -X DELETE ${AS_URL}/apps/test${i}
        dcos marathon app remove test${i}
        sleep 10s
    done
}

# Main
if ! curl ${AS_URL}; then
    echo "Please start autoscaler first."
    exit 1
fi

if  [ "${1}" == "stop" ]; then
    stop
else
    start
fi

echo "Done"