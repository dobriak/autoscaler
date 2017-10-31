FROM scratch
ADD autoscaler /
ENV AS_BASEURL=http://leader.mesos
CMD ["/autoscaler"]
