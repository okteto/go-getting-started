# Getting Started on Okteto with Go


This example shows how to use the [Okteto CLI](https://github.com/okteto/okteto) to develop a Go Sample App directly in Kubernetes. The Go Sample App is deployed using Kubernetes manifests.

This sample app is a web server that listens on port 8080 and displays workloads (*currently pod only*) in your kubernetes cluster. 
This example can only fetch workloads from the current namespace. 

This web service exposed 3 endpoints
1. '/': displays the list of workloads sorted by name
2. '/age': displays the list of workloads sorted by their creation timestamp.
3. '/restart': displays list of pods sorted by their restart count.

## Future work
* Add more unit tests especially for error cases where pods can't be fetched due to auth issues or other errors.
* Add more endpoints to fetch more information about workloads.
* Display other workloads
* Add capability to query workloads in namespaces other "default" namespace.
* Add capability to query workloads in other remote clusters.
* Make UI pretty!

