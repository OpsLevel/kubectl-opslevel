# Local development

1. [First steps | Local installation](#first-steps)
1. [Setup testing env](#setup-testing-env)
1. [Test commands locally](#test-commands-locally)
1. [Cleanup](#cleanup)

## First Steps

**Get an OpsLevel API token for a dev environment.** Do not use your prod environment since this involves mutating lots of services in OpsLevel.

Install these:

1. [minikube](https://minikube.sigs.k8s.io/docs/start/)  - this is a suggestion, you can also use **your own dev cluster if you have one.** Other options include [kind](https://kind.sigs.k8s.io/) or [k3s](https://k3s.io/).
1. [task](https://taskfile.dev/) - this will install dependencies for you.
1. [OpsLevel CLI](https://github.com/OpsLevel/cli) - makes things way easier.

```
brew install minikube
brew install task
brew install opslevel/tap/cli
```

Then run:

```
task setup
cd src                       # from now on this tutorial is going to assume you're in the src/ directory.
go build -o kubectl-opslevel .
./kubectl-opslevel           # you can also use `go run main.go` or set an alias like `alias ko='./kubectl-opslevel'`
                             # and `alias bko='go build . && ./kubectl-opslevel'`
```


> [!TIP]
> Your local version can only be called using `./kubectl-opslevel`. That is NOT THE SAME as `kubectl opslevel` which is the actual kubectl plugin you have installed. You can see what that is exactly by doing `kubectl plugin list`. You can install your local build of the plugin by copying it to /usr/local/bin and UNINSTALLING the version you have installed already (if any).

## Setup testing env

```
minikube start
```

> [!CAUTION]
> Beyond this step, make sure your kubectl is set to your minikube cluster AND that your `OPSLEVEL_API_TOKEN` is for a dev environment.

Next add some pods resources for testing and add a config file. You can use whatever resources and config you'd like, there are some examples provided in this folder.

```
kubectl apply -f ../docs/example-pod.yaml
cp ../docs/example-opslevel-k8s.yaml .
```

You can find even more examples of config files by doing `./kubectl-opslevel config sample -h`.

## Test commands locally

**Service Preview**

You can test `./kubectl-opslevel service preview` at any time by just running the command. 

Make sure the service count matches up with what you expect and that the service fields and tags are being read correctly.

```
7:01PM INF [/v1/pods] Informer is ready and synced
The following data was found in your Kubernetes cluster ...

[
    {
        "Name": "random-nginx-1",
        "Description": "stack was not necessary - rolled back",
        "Owner": "platform",
        "Aliases": [
            "default-random-nginx-1"
        ],
        "TagAssigns": [
            {
                "key": "k8s_created",
                "value": "2024-01-16T03:05:51Z"
            },
            {
                "key": "hello",
                "value": "world"
            },
```

**Service Import**

Use `kubectl edit pod/random-nginx-2` and change some fields and tags. Then import your changes with `./kubectl-opslevel service import`.

Import should loop only once and stop just like preview except it will actually update/create data in OpsLevel.

```
7:21PM INF [/v1/pods] Informer is ready and synced
7:21PM INF [random-nginx-1] No changes detected to fields - skipping update
7:21PM INF [random-nginx-1] All tags already assigned to service.
7:21PM INF [random-nginx-2] Updated Service - Diff:
  &opslevel.Service{
  	ApiDocumentPath: "",
- 	Description:     "the second nginx service",
+ 	Description:     "the **second** nginx service",
  	Framework:       "",
  	HtmlURL:         "https://app.opslevel.com/services/random-nginx-2",
  	... // 3 identical fields
  	ManagedAliases: {"default-random-nginx-2"},
  	Name:           "random-nginx-2",
  	Owner: opslevel.TeamId{
- 		Alias: "platform",
+ 		Alias: "engineering",
  		Id: strings.Join({
  			"Z2lkOi8vb3BzbGV2ZWwvVGVhbS8",
- 			"5NzU5",
+ 			"xMTA5OQ",
  		}, ""),
  	},
  	PreferredApiDocument:       nil,
  	PreferredApiDocumentSource: nil,
  	... // 2 identical fields
  	Tags: &{Nodes: {{Id: "Z2lkOi8vb3BzbGV2ZWwvVGFnLzcxMzA4ODAx", Key: "app.kubernetes.io/name", Value: "proxy"}, {Id: "Z2lkOi8vb3BzbGV2ZWwvVGFnLzcxNjMxODgz", Key: "foo", Value: "bar"}, {Id: "Z2lkOi8vb3BzbGV2ZWwvVGFnLzcxMzA5OTQ2", Key: "goodbye", Value: "world"}, {Id: "Z2lkOi8vb3BzbGV2ZWwvVGFnLzcxNjI2NjEz", Key: "hello", Value: "world"}, ...}, PageInfo: {Start: "MQ", End: "Ng"}, TotalCount: 6},
  	Tier: {},
  	Timestamps: opslevel.Timestamps{
  		CreatedAt: {Time: s"2024-01-16 03:07:03.745338 +0000 UTC"},
- 		UpdatedAt: iso8601.Time{Time: s"2024-01-18 00:20:03.530134 +0000 UTC"},
+ 		UpdatedAt: iso8601.Time{Time: s"2024-01-18 00:21:57.068903 +0000 UTC"},
  	},
  	Tools:        &{Nodes: {}},
  	Dependencies: nil,
  	... // 2 identical fields
  }
```

> [!TIP]
> If you don't see any API calls being made when importing or reconciling, it could be that your reconciler is already running in the background! Check `ps aux | grep reconcile` and try `kill -2` on those processes.

**Service Reconcile**

`./kubectl-opslevel service reconcile` will run in the background and should automatically listen for updates.

Reconcile is just like Import except it should run continuously and pickup changes in the background without exiting.

To make reconcile exit, you can send a `kill -2` (SIGINT) or `kill -15` (SIGTERM) to the process. 
SIGINT is the same as selecting the terminal in the foreground and pressing `CTRL+C`.

**Random note:** if you are using `go run main.go` the exit codes for the process will not be the same as the actual built version.

> [!TIP]
> Use `./kubectl-opslevel --log-level debug [command] [subcommand]` for more clear output on what's going on.

## Cleanup

If you're using minikube you can simply `minikube stop`.

If all your test services are named `random-nginx-DIGITS` you can use these commands to delete them easily. 

> [!CAUTION]
> Make sure you are in your dev cluster and that your current OpsLevel API token is for your dev environment.

```
echo $OPSLEVEL_API_TOKEN
opslevel list services -o json > services.json
cat services.json | jq '.[] | select(.name | match("random-nginx-\\d+")) | "\(.id)|\(.name)"' > delete.json
cat delete.json | tr -d '"' | cut -d '|' -f1 | xargs -n1 echo            # warning: double check the output with echo first
                                                                         # you can replace echo with `opslevel delete service`
```
