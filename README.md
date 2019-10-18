# Getting Started on Okteto with Golang

This example shows how to leverage [Okteto](https://github.com/okteto/okteto) to develop a Golang Sample App directly in Kubernetes. The Golang Sample App is deployed using raw Kubernetes manifests.

Okteto works in any Kubernetes cluster by reading your local kubeconfig file. If you need access to a Kubernetes cluster, [Okteto Cloud](https://cloud.okteto.com) gives you free access to secure Kubernetes namespaces, compatible with any Kubernetes tool.

## Step 1: Deploy the Go Sample App

Get a local version of the Go Sample App by executing the following commands:

```console
$ git clone https://github.com/okteto/go-getting-started
```

The `k8s.yml` file contains the raw Kubernetes manifests to deploy the Golang Sample App. Run the application by executing:

```console
$ kubectl apply -f k8s.yml
```

```console
deployment.apps "hello-world" created
service "hello-world" created
```

This is cool! You typed one command and a dev version of your application just runs 😎. 

## Step 2: Start your development environment in Kubernetes

With the Golang Sample Application deployed, run the following command:

```console
$ okteto up
 ✓  Development environment activated
 ✓  Files synchronized
    Namespace: pchico83
    Name:      hello-world
    Forward:   8080 -> 8080
               2345 -> 2345

okteto>
```

The `okteto up` command starts a [Kubernetes development environment](https://okteto.com/docs/reference/development-environment/index.html), which means:

- The Go Sample App container is updated with the docker image `okteto/golang:1`. This image contains the required dev tools to build, test and run the Go Sample App.
- A [file synchronization service](https://okteto.com/docs/reference/file-synchronization/index.html) is created to keep your changes up-to-date between your local filesystem and your application pods.
- Attach a volume to persist the Golang cache in your Kubernetes development environment.
- Container ports 8080 (the application) and 2345 (the debugger) are forwarded to localhost.
- You have a remote shell in the development environment. Build, test and run your application as if you were in your local machine.

> All of this (and more) can be customized via the `okteto.yml` [manifest file](https://okteto.com/docs/reference/manifest/index.html).


To run the application, execute in the remote shell:

```console
okteto> go run main.go
Starting hello-world server...
```

Test your application by running the command below in a local shell:

```console
$ curl localhost:8080
```

```console
Hello world!
```

## Step 3: Develop directly in Kubernetes

Opem the file `main.go` in your favorite local IDE and modify the response message on line 17 to be *Hello world from the cluster!*. Save your changes.

```golang
func helloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world from the cluster!")
}
```

Okteto will synchronize your changes to your development environnment in Kubernetes. Cancel the execution of `go run main.go` from the remote shell by pressing `ctrl + c`. Rerun your application:

```console
okteto> go run main.go
Starting hello-world server...
```

Call your application from a local shell to validate the changes:

```console
$ curl localhost:8080
```

```console
Hello world from the cluster!
```

Cool! Your code changes were instantly applied to Kubernetes. No commit, build or push required 😎!

## Step 4: Debug directly in Kubernetes

Okteto enables you to debug your applications directly from your favorite IDE. Let's take a look at how that works in VS Code, one of the most popular IDEs for Golang development.

Cancel the execution of `go run main.go` from the remote shell by pressing `ctrl + c`. Rerun your application in debug mode:

```console
okteto> dlv debug --headless --listen=:2345 --log --api-version=2
API server listening at: [::]:2345
2019-10-17T14:39:24Z info layer=debugger launching process with args: [/okteto/__debug_bin]
```

Open the _Debug_ extension and run the *Connect to okteto* launch configuration:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Connect to okteto",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "remotePath": "/okteto",
            "port": 2345,
            "host": "127.0.0.1"
        }
    ]
}
```

 Add a breakpoint on `main.go`, line 17, and call your application by running the command below from a local shell.

```console
$ curl localhost:8080
```

The execution will halt at your breakpoint. You can then inspect the request, the available variables, etc...

## Step 5: Cleanup

Cancel the `okteto up` command by pressing `ctrl + c` + `ctrl + d` and run the following commands to remove the resources created by this guide: 

```console
$ okteto down
 ✓  Development environment deactivated
```

```console
$ kubectl delete -f k8s.yml
deployment.apps "hello-world" deleted
service "hello-world" deleted
```
