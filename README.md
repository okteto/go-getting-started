# Develop Go Applications directly in Kubernetes with Okteto

This example shows how to leverage [Okteto](https://github.com/okteto/okteto) to develop a Go Sample App directly in Kubernetes. 

Okteto works in any Kubernetes cluster by reading your local Kubernetes credentials. For a more empowered experience, follow this [guide](https://okteto.com/docs/samples/golang/) to deploy the Go Sample App in [Okteto Cloud](https://cloud.okteto.com), a free Kubernetes cluster for developers.

## Step 1: Install the Okteto CLI

Install the Okteto CLI by following our [installation guides](https://github.com/okteto/okteto/blob/master/docs/installation.md).

## Step 2: Deploy the Go Sample App

Get a local version of the Go Sample App by executing the following commands in your local terminal:

```console
$ git clone https://github.com/okteto/go-getting-started
```

In the `manifest/` directory you have the Kubernetes manifests that we will use in this guide to deploy the application in the cluster. Okteto works with any Kubernetes tool or deployment practice.

> If you don't have `kubectl` installed, follow this [guide](https://kubernetes.io/docs/tasks/tools/install-kubectl/).

Deploy the Go Sample App by executing:

```console
$ kubectl apply -f manifests
deployment.apps "hello-world" created
service "hello-world" created
```

## Step 3: Create your Okteto Environment

With the app deployed, you can start your Okteto Environment by running the following command:

```console
$ okteto up
 ✓  Okteto Environment activated
 ✓  Files synchronized
 ✓  Your Okteto Environment is ready
    Namespace: cindy
    Name:      hello-world
    Forward:   8080 -> 8080

okteto>
```

The `okteto up` command will automatically start an Okteto Environment, which means:

- The Go Sample App container is updated with the docker image `okteto/golang:1`. This image contains the required dev tools to build, test and run the Go Sample App.
- A bidirectional file synchronization service is started to keep your changes up to date between your local filesystem and your Okteto Environment.

Once the Okteto Environment is ready, start your application by executing the following command in your Okteto Terminal:

```console
okteto> go run main.go
Starting hello-world server...
```

You can now access the Go Sample App at http://localhost:8080.

## Step 4: Develop directly in the cloud

Now things get more exciting. Edit the file `main.go` and replace the word `cluster` with `Okteto` on line 23. Save your changes.

Cancel the execution of `go run main.go` from your Okteto Terminal by pressing `ctrl + c`. Now rerun your application:

```console
okteto> go run main.go
```

Go back to the browser and reload the page. Notice how your changes are instantly applied. No commit, build or push required 😎! 


## Step 5: Cleanup

Cancel the `okteto up` command by pressing `ctrl + c` + `exit` and run the following commands to remove the resources created by this guide: 

```console
$ okteto down -v
 ✓  Okteto Environment deactivated
 
```

```console
$ kubectl delete -f manifests
deployment.apps "hello-world" deleted
service "hello-world" deleted
```
