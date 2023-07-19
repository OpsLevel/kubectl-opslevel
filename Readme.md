<p align="center">
    <a href="https://github.com/OpsLevel/kubectl-opslevel/blob/main/LICENSE" alt="License">
        <img src="https://img.shields.io/github/license/OpsLevel/kubectl-opslevel.svg" /></a>
    <a href="https://pkg.go.dev/github.com/OpsLevel/kubectl-opslevel">
        <img src="https://pkg.go.dev/badge/github.com/OpsLevel/kubectl-opslevel" alt="GoDoc"></a>
    <a href="https://goreportcard.com/report/github.com/OpsLevel/kubectl-opslevel" alt="Go Report Card">
        <img src="https://goreportcard.com/badge/github.com/OpsLevel/kubectl-opslevel" /></a>
    <a href="https://GitHub.com/OpsLevel/kubectl-opslevel/releases/" alt="Release">
        <img src="https://img.shields.io/github/v/release/OpsLevel/kubectl-opslevel" /></a> 
    <a href="https://masterminds.github.io/stability/active.html" alt="Stability: Active">
        <img src="https://masterminds.github.io/stability/active.svg" /></a>  
    <a href="https://github.com/OpsLevel/kubectl-opslevel/graphs/contributors" alt="Contributors">
        <img src="https://img.shields.io/github/contributors/OpsLevel/kubectl-opslevel" /></a>
    <a href="https://github.com/OpsLevel/kubectl-opslevel/pulse" alt="Activity">
        <img src="https://img.shields.io/github/commit-activity/m/OpsLevel/kubectl-opslevel" /></a>
    <a href="https://github.com/OpsLevel/kubectl-opslevel/releases" alt="Downloads">
        <img src="https://img.shields.io/github/downloads/OpsLevel/kubectl-opslevel/total" /></a>
</p>

<p align="center">
 <a href="#prerequisite">Prerequisite</a> |
 <a href="#installation">Installation</a> |
 <a href="#quickstart">Quickstart</a> |
 <a href="https://www.opslevel.com/docs/integrations/kubernetes/">Documentation</a> |
 <a href="#troubleshooting">Troubleshooting</a>
</p>

[![Overall](https://img.shields.io/endpoint?style=flat&url=https%3A%2F%2Fapp.opslevel.com%2Fapi%2Fservice_level%2F4SZo_XBzNM8K84zLHYdEXCcvBL6q_pTzUMSR09DmnZM)](https://app.opslevel.com/services/opslevel_kubernetes_sync/maturity-report)

`kubectl-opslevel` is a command line tool that enables you to import & reconcile services with [OpsLevel](https://www.opslevel.com/) from your Kubernetes clusters.  You can also run this tool inside your Kubernetes cluster as a job to reconcile the data with OpsLevel periodically using our [Helm Chart](https://github.com/OpsLevel/helm-charts).

## Prerequisite

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [jq](https://stedolan.github.io/jq/download/)
- [OpsLevel API Token](https://app.opslevel.com/api_tokens)

## Installation

```sh
brew install opslevel/tap/kubectl
```

<!--
#### Deb

```sh
sudo apt-get install apt-transport-https
wget -qO - https://opslevel.github.io/kubectl-opslevel-repo/deb/public.key | sudo apt-key add -
echo deb https://opslevel.github.io/kubectl-opslevel-repo/deb [CODE_NAME] main | sudo tee -a /etc/apt/sources.list
sudo apt-get update
sudo apt-get install kubectl-opslevel
```

#### RPM

```sh
cat << EOF > /etc/yum.repos.d/opslevel.repo
[opslevel]
name=opslevel cli repository
baseurl=https://opslevel.github.io/kubectl-opslevel-repo/rpm/releases/$releasever/$basearch/
gpgcheck=0
enabled=1
EOF
sudo yum -y update
sudo yum -y install kubectl-opslevel
```
-->

## Docker

The docker container is hosted on [AWS Public ECR](https://gallery.ecr.aws/opslevel/kubectl-opslevel)

## Quickstart

```sh
# Generate a config file
kubectl opslevel config sample > ./opslevel-k8s.yaml

# Like Terraform, generate a preview of data from your Kubernetes cluster
# NOTE: this step does not validate any of the data with OpsLevel
kubectl opslevel service preview

# Import (and reconcile) the found data with your OpsLevel account
 OPSLEVEL_API_TOKEN=XXXX kubectl opslevel service import
```

[![asciicast](https://asciinema.org/a/bv6WTcqkGtmC5wXN4VXYr035y.svg)](https://asciinema.org/a/bv6WTcqkGtmC5wXN4VXYr035y)


### Current Sample Configuration

```yaml
version: "1.1.0"
service:
  import:
    - selector: # This limits what data we look at in Kubernetes
        apiVersion: apps/v1 # only supports resources found in 'kubectl api-resources --verbs="get,list"'
        kind: Deployment
        excludes: # filters out resources if any expression returns truthy
          - .metadata.namespace == "kube-system"
          - .metadata.annotations."opslevel.com/ignore"
      opslevel: # This is how you map your kubernetes data to opslevel service
        name: .metadata.name
        description: .metadata.annotations."opslevel.com/description"
        owner: .metadata.annotations."opslevel.com/owner"
        lifecycle: .metadata.annotations."opslevel.com/lifecycle"
        tier: .metadata.annotations."opslevel.com/tier"
        product: .metadata.annotations."opslevel.com/product"
        language: .metadata.annotations."opslevel.com/language"
        framework: .metadata.annotations."opslevel.com/framework"
        system: .metadata.annotations."opslevel.com/system"
        aliases: # This are how we identify the services again during reconciliation - please make sure they are very unique
          - '"k8s:\(.metadata.name)-\(.metadata.namespace)"'
        tags:
          assign: # tag with the same key name but with a different value will be updated on the service
            - '{"imported": "kubectl-opslevel"}'
            # find annoations with format: opslevel.com/tags.<key name>: <value>
            - '.metadata.annotations | to_entries |  map(select(.key | startswith("opslevel.com/tags"))) | map({(.key | split(".")[2]): .value})'
            - .metadata.labels
          create: # tag with the same key name but with a different value with be added to the service
            - '{"environment": .spec.template.metadata.labels.environment}'
        tools:
          - '{"category": "other", "displayName": "my-cool-tool", "url": .metadata.annotations."example.com/my-cool-tool"} | if .url then . else empty end'
          # find annotations with format: opslevel.com/tools.<category>.<displayname>: <url> 
          - '.metadata.annotations | to_entries |  map(select(.key | startswith("opslevel.com/tools"))) | map({"category": .key | split(".")[2], "displayName": .key | split(".")[3], "url": .value})'
        repositories: # attach repositories to the service using the opslevel repo alias - IE github.com:hashicorp/vault
          - '{"name": "My Cool Repo", "directory": "/", "repo": .metadata.annotations.repo} | if .repo then . else empty end'
          # if just the alias is returned as a single string we'll build the name for you and set the directory to "/"
          - .metadata.annotations.repo
          # find annotations with format: opslevel.com/repo.<displayname>.<repo.subpath.dots.turned.to.forwardslash>: <opslevel repo alias> 
          - '.metadata.annotations | to_entries |  map(select(.key | startswith("opslevel.com/repos"))) | map({"name": .key | split(".")[2], "directory": .key | split(".")[3:] | join("/"), "repo": .value})'
```

### Enable shell autocompletion

We have the ability to generate autocompletion scripts for the shell's `bash`, `zsh`, `fish` and `powershell`.  To generate 
the completion script for MacOS zsh:

```sh
kubectl opslevel completion zsh > /usr/local/share/zsh/site-functions/_kubectl-opslevel
```

Make sure you have `zsh` completion turned on by having the following as one of the first few lines in your `.zshrc` file

```sh
echo "autoload -U compinit; compinit" >> ~/.zshrc
```

### Collect

This tool also comes with a way to collect the Kubernetes data to our custom event checks (CEC) as payloads.

```yaml
#Sample Opslevel CLI Config
version: "1.2.0"
service:
  import:
    ... Import Coammands Config Here ...
  collect:
    - selector: # This limits what data we look at in Kubernetes
        apiVersion: apps/v1 # only supports resources found in 'kubectl api-resources --verbs="get,list"'
        kind: Deployment
        excludes: # filters out resources if any expression returns truthy
          - .metadata.namespace == "kube-system"
          - .metadata.annotations."opslevel.com/ignore"
```

```sh
kubectl opslevel collect --integration-url https://app.opslevel.com/integrations/custom_event/XXXX-XX-XX-XXXXX
```

This will then use the `collect` stanza's selector to get the JSON payload of the resource and ship it to our CEC endpoint for use in the maturity rubric checks.

### JSON-Schema

The tool also has the ability to output a [JSON-Schema](https://json-schema.org/) file for use in IDE's when editing the configuration file.
You can read more about adding JSON-Schema validate to [VS Code](https://code.visualstudio.com/docs/languages/json#_json-schemas-and-settings)

```sh
kubectl opslevel config schema > ~/.opslevel-k8s-schema.json
```

Then add the following to you [VS Code user settings](https://code.visualstudio.com/docs/getstarted/settings)

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

Generally speaking if any other command works IE `kubectl get deployment` then any `kubectl opslevel` command should work too.  If this is the not the case then there is likely a special authentication mechanism in place that we are not handling properly.  This should be reported as a bug.

### A field mapped in the configuration file is not in the service data

For the most part `jq` filter failures are bubbled up but in certain edgecases they can fail silently.
The best way to test a `jq` expression in isoloation is to emit the Kubernetes resource to json IE `kubectl get deployment <name> -o json` 
and then play around with the expression in [jqplay](https://jqplay.org/)

Generally speaking if we detect a json `null` value we do build any data for that field.

### String interpolation has NULL in it

There is a special edgecase with string interpolation and null values that we cannot handle that is documented [here](/../../issues/36)

### Unable to list all Namespaces

Sometimes in tight permissions cluster listing of all Namespaces is not allowed.  The tool currently tries to list all Namespaces
in a cluster to use as a batching mechanism.  This functionality can be skipped by using
an the explict list `namespaces` in the selector which skips the API call to Kubernetes to list all Namespaces.

```yaml
service:
  import:
    - selector: # This limits what data we look at in Kubernetes
        apiVersion: apps/v1 # only supports resources found in 'kubectl api-resources --verbs="get,list"'
        kind: Deployment
        namespaces:
          - default
          - kube-system
```
