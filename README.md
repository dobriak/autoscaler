# autoscaler

## DC/OS EE autoscaler for Marathon apps

Multithreaded Marathon app metrics monitor with RESTful interface.
Manages Marathon app moniroting threads that can issue scaling signals based on a combination of metrics and conditions.

### RESTful APIs:

* GET /apps - returns a list of running Marathon apps scaling monitors
* GET /app  -d '{"app_id":"/some/app"}' - returns a single Marathon app scaling monitor
* POST /apps -d '{"app_id": "/some/app", "max_cpu": 45 ...}' - create and start a Marathon app scaling monitor. See below for all options.
* DELETE /apps -d '{"app_id":"/some/app"}' - stop and remove a single Marathon app sclaing monitor

### Examples

Get a list of all Marathon app scaling monitors:
``` bash
curl http://autoscaler.marathon.l4lb.thisdcos.directory/apps
```

Retrieve Marathon app scaling monitor with ID of ```test1```:
``` bash
curl -d '{"app_id":"/test1"}' http://autoscaler.marathon.l4lb.thisdcos.directory/app
```

Create and start a Marathon app scaling monitor wit ID of ```/myapp/test1```:
```bash
curl -X POST -d '{ "app_id": "/myapp/test1","max_cpu": 60,"min_cpu": 20,"max_mem": 90,"min_mem": 5,"method": "cpu","scale_factor": 1,"max_instances": 6,"min_instances": 2,"warm_up": 3,"cool_down": 3,"interval": 181}' http://autoscaler.marathon.l4lb.thisdcos.directory/apps
```

Stop and remove Marathon app scaling monitor with ID of ```test3```:
```bash
curl -X DELETE -d '{"app_id":"/test3"}' http://autoscaler.marathon.l4lb.thisdcos.directory/apps
```

### App scaling monitor options
* **app_id** - App monitor ID, same as the Marathon app ID being monitored
* **max_cpu** - Maximum percentage of CPU utilization inside the container allowed
* **min_cpu** - Minimum percentage of CPU
* **max_mem** - Maximum percentage of available to container memory utilization
* **min_mem** - Minimum percentage of memory
* **method** - Method to use for signal generation. It can be: memory, CPU, and, or (mem|cpu|and|or)
* **scale_factor** - Instances number multiplier when scaling up or down
* **max_instances** - Maximum number of instances per monitored app
* **min_instances** - Minimum number of instances
* **warm_up** - How many monitoring cycles with "up" signal before scaling up
* **cool_down** - How many monitoring cycles with "down" signal before scaling down
* **interval** - Length of monitoring cycle in seconds

### Usage
This application is meant to be run as a Marathon service on your Enterprise Edition DC/OS cluster version 1.10 and up. In order for it to be able to scale Marathon services it will need a set of specific permissions and a service account. 

This repository includes a script that makes it easy to create a service account and assign it to a namespace. Run the following from a machine that has the DC/OS CLI installed and configured:

``` bash
./create-service-account.sh scaler myapp
```

Here ```scaler``` is the name of the service account and ```myapp``` is the namespace (folder) in which you can spin up Marathon apps and create autoscaling app monitors. The script creates a secret, the contents of which will have to be assigned to an environment variable inside the ```autoscaler``` container.

The application itself is Dockerized and available in DockerHub under ```dobriak/autoscaler:0.0.1```. It offers its RESTful interface on port ```8080``` container-side.

Once the service account is created, create a Marathon app from the included ```marathon-autoscaler.json``` file. From the same DC/OS CLI equipped machine, run

``` bash
dcos marathon app add marathon-autoscaler.json
```

The app definition itself sets up an internal named VIP that exposes said RESTful API at ```http://autoscaler.marathon.l4lb.thisdcos.directory```. It also assigns the value of the automatically created secret to the ```AS_SECRET``` environment variable needed for autoscaler to work.

### Compiling and building the Docker container yourself

The Dockefile is included. One thing to note is since the base container is ```scratch``` you will have to do a static compilation before doing ```docker build```

``` bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o autoscaler .
```

### Testing the autoscaler

Once the autoscaler is installed and in a healthy state you can quickly spin up a few Marathon apps and create monitors for them. To make things easier, I've included a script with 4 sample apps. Run this script from _inside_ the cluster:

``` bash
[laptop] scp -r -i /path/to/key sample_apps.sh test/ <user>@<master or node ip>:
# Make sure to edit the top few lines of the script, for example
# provide your super username and password
[inside-the-cluster] ./sample_apps.sh
```

To clean up Marathon apps and app monitors:

``` bash
[inside-the-cluster] ./sample_apps.sh stop
```

That's it!
Issues and PRs are welcome!


