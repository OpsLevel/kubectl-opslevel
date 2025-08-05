> [!IMPORTANT]  
> This project and tool as been deprecated in favor of the [opslevel-agent](https://github.com/OpsLevel/helm-charts/blob/main/charts/opslevel-agent/README.md#opslevel-agent).  Please Migrate!



<p align="center">
    <a href="https://github.com/OpsLevel/kubectl-opslevel/blob/main/LICENSE">
        <img src="https://img.shields.io/github/license/OpsLevel/kubectl-opslevel.svg" alt="License" /></a>
    <a href="https://pkg.go.dev/github.com/OpsLevel/kubectl-opslevel">
        <img src="https://pkg.go.dev/badge/github.com/OpsLevel/kubectl-opslevel" alt="GoDoc" /></a>
    <a href="https://goreportcard.com/report/github.com/OpsLevel/kubectl-opslevel">
        <img src="https://goreportcard.com/badge/github.com/OpsLevel/kubectl-opslevel" alt="Go Report Card" /></a>
    <a href="https://GitHub.com/OpsLevel/kubectl-opslevel/releases/">
        <img src="https://img.shields.io/github/v/release/OpsLevel/kubectl-opslevel" alt="Release" /></a>
    <a href="https://masterminds.github.io/stability/active.html">
        <img src="https://masterminds.github.io/stability/active.svg" alt="Stability: Active" /></a>
    <a href="https://github.com/OpsLevel/kubectl-opslevel/graphs/contributors">
        <img src="https://img.shields.io/github/contributors/OpsLevel/kubectl-opslevel" alt="Contributors" /></a>
    <a href="https://github.com/OpsLevel/kubectl-opslevel/pulse">
        <img src="https://img.shields.io/github/commit-activity/m/OpsLevel/kubectl-opslevel" alt="Activity" /></a>
    <a href="https://github.com/OpsLevel/kubectl-opslevel/releases">
        <img src="https://img.shields.io/github/downloads/OpsLevel/kubectl-opslevel/total" alt="Downloads" /></a>
</p>

<p align="center">
 <a href="#prerequisite">Prerequisite</a> |
 <a href="#installation">Installation</a> |
 <a href="./CONTRIBUTING.md">CONTRIBUTING</a> |
 <a href="./docs/Local-Development.md">Local Dev/Testing Instructions</a> |
 <a href="#quickstart">Quickstart</a> |
 <a href="https://docs.opslevel.com/docs/kubernetes-integration">Documentation</a> |
 <a href="#troubleshooting">Troubleshooting</a>
</p>

[![Overall](https://img.shields.io/endpoint?style=flat&url=https%3A%2F%2Fapp.opslevel.com%2Fapi%2Fservice_level%2F4SZo_XBzNM8K84zLHYdEXCcvBL6q_pTzUMSR09DmnZM)](https://app.opslevel.com/services/opslevel_kubernetes_sync/maturity-report)

`kubectl-opslevel` is a command line tool that enables you to import & reconcile services with [OpsLevel](https://www.opslevel.com/) from your Kubernetes clusters.  You can also run this tool inside your Kubernetes cluster as a job to reconcile the data with OpsLevel periodically using our [Helm Chart](https://github.com/OpsLevel/helm-charts).

## Prerequisites

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [jq](https://jqlang.github.io/jq/download/)
- [OpsLevel API Token](https://app.opslevel.com/api_tokens)

## Installation

```sh
# OR - manually copy the binary to /usr/bin/local
brew install opslevel/tap/kubectl
```

## Docker

The docker container is hosted on [AWS Public ECR](https://gallery.ecr.aws/opslevel/kubectl-opslevel)

## Quickstart

```sh
# Generate a config file
kubectl opslevel config sample > ./opslevel-k8s.yaml

# Like Terraform, generate a preview of data from your Kubernetes cluster
# NOTE: this step does not validate any of the data with OpsLevel
 OPSLEVEL_API_TOKEN=XXXX kubectl opslevel service preview

# Import (and reconcile) the found data with your OpsLevel account
 OPSLEVEL_API_TOKEN=XXXX kubectl opslevel service import
```

[![asciicast](https://asciinema.org/a/bv6WTcqkGtmC5wXN4VXYr035y.svg)](https://asciinema.org/a/bv6WTcqkGtmC5wXN4VXYr035y)


### Current Sample Configuration

```yaml
version: "1.3.0"
service:
  import:
    - selector:
        apiVersion: apps/v1
        kind: Deployment
        excludes:
          - .metadata.namespace == "kube-system"
          - .metadata.annotations."opslevel.com/ignore"
      opslevel:
        aliases:
          - '"k8s:\(.metadata.name)-\(.metadata.namespace)"'
          - '"\(.metadata.namespace)-\(.metadata.name)"'
        description: .metadata.annotations."opslevel.com/description"
        framework: .metadata.annotations."opslevel.com/framework"
        language: .metadata.annotations."opslevel.com/language"
        lifecycle: .metadata.annotations."opslevel.com/lifecycle"
        name: .metadata.name
        owner: .metadata.annotations."opslevel.com/owner"
        product: .metadata.annotations."opslevel.com/product"
        properties:
          prop_object: .metadata.annotations.prop_object
        repositories:
          - '{"name": "My Cool Repo", "directory": "", "repo": .metadata.annotations.repo} | if .repo then . else empty end'
          - .metadata.annotations.repo
          - '.metadata.annotations | to_entries |  map(select(.key | startswith("opslevel.com/repo"))) | map({"name": .key | split(".")[2], "directory": .key | split(".")[3:] | join("/"), "repo": .value})'
        system: ""
        tags:
          assign:
            - '{"imported": "kubectl-opslevel"}'
            - '.metadata.annotations | to_entries |  map(select(.key | startswith("opslevel.com/tags"))) | map({(.key | split(".")[2]): .value})'
            - .metadata.labels
          create:
            - '{"environment": .spec.template.metadata.labels.environment}'
        tier: .metadata.annotations."opslevel.com/tier"
        tools:
          - '{"category": "other", "environment": "production", "displayName": "my-cool-tool", "url": .metadata.annotations."example.com/my-cool-tool"} | if .url then . else empty end'
          - '.metadata.annotations | to_entries |  map(select(.key | startswith("opslevel.com/tools"))) | map({"category": .key | split(".")[2], "displayName": .key | split(".")[3], "url": .value})'
```

### Enable shell autocompletion

We have the ability to generate autocompletion scripts for the shell's `bash`, `zsh`, `fish` and `powershell`.  To generate
the completion script for macOS zsh:

```sh
kubectl opslevel completion zsh > /usr/local/share/zsh/site-functions/_kubectl-opslevel
```

Make sure you have `zsh` completion turned on by having the following as one of the first few lines in your `.zshrc` file

```sh
echo "autoload -U compinit; compinit" >> ~/.zshrc
```

### JSON-Schema

The tool also has the ability to output a [JSON-Schema](https://json-schema.org/) file for use in IDEs when editing the configuration file.
You can read more about adding JSON-Schema validate to [VS Code](https://code.visualstudio.com/docs/languages/json#_json-schemas-and-settings)

```sh
kubectl opslevel config schema > ~/.opslevel-k8s-schema.json
```

Then add the following to your [VS Code user settings](https://code.visualstudio.com/docs/getstarted/settings)

```json
    "yaml.schemas": {
        "~/.opslevel-k8s-schema.json": ["opslevel-k8s.yaml"],
    }
```

## Troubleshooting

### No services output from `service preview`

This can happen for a number of reasons:

  - Kubernetes RBAC permissions do not allow for listing namespaces
  - Configuration file exclude rules exclude all found resources

### Unable to connect to Kubernetes cluster

Generally speaking if any other command works IE `kubectl get deployment` then any `kubectl opslevel` command should work too.  If this is not the case then there is likely a special authentication mechanism in place that we are not handling properly.  This should be reported as a bug.

### A field mapped in the configuration file is not in the service data

For the most part `jq` filter failures are bubbled up but in certain edgecases they can fail silently.
The best way to test a `jq` expression in isoloation is to emit the Kubernetes resource to json IE `kubectl get deployment <name> -o json`
and then play around with the expression in [jqplay](https://jqplay.org/)

Generally speaking if we detect a json `null` value we don't build any data for that field.

### String interpolation has NULL in it

There is a special edgecase with string interpolation and null values that we cannot handle that is documented [here](https://github.com/OpsLevel/kubectl-opslevel/issues/36)

### Unable to list all Namespaces

Sometimes in tight permissions cluster listing of all Namespaces is not allowed.  The tool currently tries to list all Namespaces
in a cluster to use as a batching mechanism.  This functionality can be skipped by using
the explict list `namespaces` in the selector which skips the API call to Kubernetes to list all Namespaces.

```yaml
version: "1.3.0"
service:
  import:
    - selector: # This limits what data we look at in Kubernetes
        apiVersion: apps/v1 # only supports resources found in 'kubectl api-resources --verbs="get,list"'
        kind: Deployment
        namespaces:
          - default
          - kube-system
```
