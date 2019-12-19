# Getting Started on Okteto with Go

This example shows how to leverage [Okteto](https://github.com/okteto/okteto) to develop a Go Sample App directly in Kubernetes. The Go Sample App is deployed using raw Kubernetes manifests.

Okteto is a client-side only tool that works in any Kubernetes cluster. If you need access to a Kubernetes cluster, [Okteto Cloud](https://cloud.okteto.com) gives you free access to sandboxes Kubernetes namespaces, compatible with any Kubernetes tool.

## Step 1: Deploy the Go Sample App

Get a local version of the Go Sample App by executing the following commands:

```console
$ git clone https://github.com/okteto/go-getting-started
```

The `k8s.yml` file contains the raw Kubernetes manifests to deploy the Go Sample App. Run the application by executing:

> If you don't have `kubectl` installed, follow this [guide](https://kubernetes.io/docs/tasks/tools/install-kubectl/).

```console
$ kubectl apply -f k8s.yml
```

```console
deployment.apps "hello-world" created
service "hello-world" created
```

This is cool! You typed one command and a dev version of your application just runs 😎. 

## Step 2: Start your development environment in Kubernetes

With the Go Sample Application deployed, run the following command:

```console
$ okteto up
 ✓  Development environment activated
 ✓  Files synchronized
    Namespace: pchico83
    Name:      hello-world
    Forward:   8080 -> 8080
               2345 -> 2345

Welcome to your development environment. Happy coding!
okteto>
```

The `okteto up` command starts a [Kubernetes development environment](https://okteto.com/docs/reference/development-environment/index.html), which means:

- The Go Sample App container is updated with the docker image `okteto/golang:1`. This image contains the required dev tools to build, test and run the Go Sample App.
- A [file synchronization service](https://okteto.com/docs/reference/file-synchronization/index.html) is created to keep your changes up-to-date between your local filesystem and your application pods.
- A volume is attached to persist the Go cache and packages in your Kubernetes development environment.
- Container ports 8080 (the application) and 2345 (the debugger) are forwarded to localhost.
- A remote shell is started in your Kubernetes development environment. Build, test and run your application as if you were in your local machine.

> All of this (and more) can be customized via the `okteto.yml` [manifest file](https://okteto.com/docs/reference/manifest/index.html). You can also use the file `.stignore` to skip files from file synchronization. This is useful to avoid synchronizing binaries or git metadata.

To run the application, execute in the remote shell:

```console
okteto> go run main.go
```

```console
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

Open the `main.go` file in your favorite local IDE and modify the response message on line 17 to be *Hello world from the cluster!*. Save your changes.

```golang
func helloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world from the cluster!")
}
```

Okteto will synchronize your changes to your development environment in Kubernetes. Cancel the execution of `go run main.go` from the remote shell by pressing `ctrl + c`. Rerun your application:

```console
okteto> go run main.go
```

```console
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
```

```console
API server listening at: [::]:2345
2019-10-17T14:39:24Z info layer=debugger launching process with args: [/okteto/__debug_bin]
```

Open the _Debug_ extension and run the *Connect to okteto* debug configuration (or press the F5 shortcut):

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

![Debug directly in Kubernetes](images/halt.png)

## Step 5: Cleanup

Cancel the `okteto up` command by pressing `ctrl + c` + `ctrl + d` and run the following commands to remove the resources created by this guide: 

```console
$ okteto down
```

```console
 ✓  Development environment deactivated
```

```console
$ kubectl delete -f k8s.yml
```

```console
deployment.apps "hello-world" deleted
service "hello-world" deleted
```
