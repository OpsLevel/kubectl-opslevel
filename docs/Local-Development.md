# Local development

1. [First steps | Local installation](#first-steps)
1. [Setup testing env](#setup-testing-env)
1. [Test commands locally](#test-commands-locally)

## First Steps

1. Install [task](https://taskfile.dev/) to setup dependencies.
1. `cd` into the repo and run `task setup`
1. Build a binary by running `cd src` and `go build -o kubectl-opslevel .`

**Notice: `./kubectl-opslevel` is not the same as `kubectl opslevel` that you may have installed from brew.**

## Setup testing env

1. **Get an [OpsLevel API Token](https://app.opslevel.com/api_tokens) for a dev (not prod!) environment.**
1. Get access to a cluster to import deployments from.

If you don't have a cluster, you can use [kind](https://kind.sigs.k8s.io/)
or [minikube](https://minikube.sigs.k8s.io/docs/) to create one.

You should be able to `kubectl get deploy` in the cluster. If there are no deployments,
you can `kubectl apply -f` our [sample deployment](./deployment.yaml) (make sure to clean it up after.)

## Test commands locally

```
# test basic functionality
# optional: add '--log-level debug'
./kubectl-opslevel version
./kubectl-opslevel config sample > opslevel-k8s.yaml
./kubectl-opslevel config view

# test detection and create/update services
# REMINDER: use a dev token, not a prod token!
# optional: add '--disable-service-create'
 OPSLEVEL_API_TOKEN=XXXX ./kubectl-opslevel service preview
 OPSLEVEL_API_TOKEN=XXXX ./kubectl-opslevel service import       # runs loop once
 OPSLEVEL_API_TOKEN=XXXX ./kubectl-opslevel service reconcile    # runs continuously

# notice: if you are using 'go run main.go' the exit codes will not be the same as the binary.
ps aux | grep -i 'kubectl-opslevel'
kill -2 <PID>                           # SIGINT - the program should exit gracefully, unlike with SIGTERM.
kill -15 <PID>                          # SIGTERM
```
