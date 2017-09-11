# autoscaler

## DC/OS EE autoscaler for Marathon

Manages Marathon moniroting threads that can issue scaling signals based on a combination of metrics and conditions.

RESTful APIs:

* GET /apps - returns a list of running Marathon apps scaling monitors
* GET /apps/{app-id} - returns a single Marathon app scaling monitor
* POST /apps - create and start a Marathon app scaling monitor
* DELETE /apps/{app-id} - stop and remove a single Marathon app sclaing monitor

### Examples

Get a list of all Marathon app scaling monitors:
```bash
curl http://localhost:8080/apps
```

Retrieve Marathon app scaling monitor with ID of ```test1```:
```bash
curl http://localhost:8080/apps/test1
```

Create and start a Marathon app scaling monitor wit ID of ```test3```:
```bash
curl -X POST -d '{ "app_id": "/test1","max_cpu": 60,"min_cpu": 20,"max_mem": 90,"min_mem": 5,"method": "cpu","scale_factor": 1,"max_instances": 6,"min_instances": 2,"warm_up": 3,"cool_down": 3,"interval": 17}' http://localhost:8080/apps
```

Stop and remove Marathon app scaling monitor with ID of ```test3```:
```bash
curl -X DELETE http://localhost:8080/apps/test3
```


