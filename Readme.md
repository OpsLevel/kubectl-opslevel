<p align="center">
    <a href="https://github.com/OpsLevel/kubectl-opslevel/blob/main/LICENSE" alt="License">
        <img src="https://img.shields.io/github/license/OpsLevel/kubectl-opslevel.svg" /></a>
    <a href="http://golang.org" alt="Made With Go">
        <img src="https://img.shields.io/github/go-mod/go-version/OpsLevel/kubectl-opslevel?filename=src%2Fgo.mod" /></a>
    <a href="https://GitHub.com/OpsLevel/kubectl-opslevel/releases/" alt="Release">
        <img src="https://img.shields.io/github/v/release/OpsLevel/kubectl-opslevel" /></a>  
    <a href="https://GitHub.com/OpsLevel/kubectl-opslevel/issues/" alt="Issues">
        <img src="https://img.shields.io/github/issues/OpsLevel/kubectl-opslevel.svg" /></a>  
    <a href="https://github.com/OpsLevel/kubectl-opslevel/graphs/contributors" alt="Contributors">
        <img src="https://img.shields.io/github/contributors/OpsLevel/kubectl-opslevel" /></a>
    <a href="https://github.com/OpsLevel/kubectl-opslevel/pulse" alt="Activity">
        <img src="https://img.shields.io/github/commit-activity/m/OpsLevel/kubectl-opslevel" /></a>
    <a href="https://dependabot.com/" alt="Dependabot">
        <img src="https://badgen.net/badge/Dependabot/enabled/green?icon=dependabot" /></a>
</p>

`kubectl-opslevel` is a command line tool that enables you to import & reconcile services with [OpsLevel](https://www.opslevel.com/) from your Kubernetes clusters.  You can also run this tool inside your Kubernetes cluster as a job to reconcile the data with OpsLevel periodically.  If you opt for this please read our [service aliases](#aliases) section as we use these to properly find and reconcile the data so it is important you choose something unique.

### Quickstart

Follow the [installation](#installation) instructions before running the below commands

```bash
# Generate a config file
kubectl opslevel config sample > ./opslevel-k8s.yaml

# Like Terraform, generate a preview of data from your Kubernetes cluster
# NOTE: this step does not validate any of the data with OpsLevel
kubectl opslevel service preview

# Import (and reconcile) the found data with your OpsLevel account
 OL_APITOKEN=XXXX kubectl opslevel service import
```

[![asciicast](https://asciinema.org/a/bv6WTcqkGtmC5wXN4VXYr035y.svg)](https://asciinema.org/a/bv6WTcqkGtmC5wXN4VXYr035y)


<blockquote>This tool is still in beta.  It's sufficently stable for production use and has successfully imported & reconciled multiple OpsLevel accounts</blockquote>

#### Current Sample Configuration

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

Table of Contents
=================

   * [Prerequisite](#prerequisite)
   * [Installation](#installation)
      * [Validate Install](#validate-install)
   * [Documentation](#documentation)

## Prerequisite

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [jq](https://stedolan.github.io/jq/download/)
- [OpsLevel API Token](https://app.opslevel.com/api_tokens)

## Installation

#### MacOS

```sh
TOOL_VERSION=$(curl -s https://api.github.com/repos/opslevel/kubectl-opslevel/releases/latest | grep tag_name | cut -d '"' -f 4)
curl -Lo kubectl-opslevel.tar.gz https://github.com/opslevel/kubectl-opslevel/releases/download/${TOOL_VERSION}/kubectl-opslevel-darwin-amd64.tar.gz
tar -xzvf kubectl-opslevel.tar.gz  
rm kubectl-opslevel.tar.gz
sudo mv kubectl-opslevel /usr/local/bin/kubectl-opslevel
```

#### Linux

```sh
TOOL_VERSION=$(curl -s https://api.github.com/repos/opslevel/kubectl-opslevel/releases/latest | grep tag_name | cut -d '"' -f 4)
curl -Lo kubectl-opslevel https://github.com/opslevel/kubectl-opslevel/releases/download/${TOOL_VERSION}/kubectl-opslevel-linux-amd64.tar.gz
tar -xzvf kubectl-opslevel.tar.gz  
rm kubectl-opslevel.tar.gz
sudo mv kubectl-opslevel /usr/local/bin/kubectl-opslevel
```

#### Docker

The docker container is hosted on [AWS Public ECR](https://gallery.ecr.aws/e1n4f2i6/kubectl-opslevel)

The following downloads the container and creates a shim at `/usr/local/bin` so you can use the binary like its downloaded natively - it will just be running in a docker container. *NOTE: you may need to adjust how your kube config is mounted and set inside the container*

```
TOOL_VERSION=$(curl -s https://api.github.com/repos/opslevel/kubectl-opslevel/releases/latest | grep tag_name | cut -d '"' -f 4)
docker pull public.ecr.aws/opslevel/kubectl-opslevel:${TOOL_VERSION}
docker tag public.ecr.aws/opslevel/kubectl-opslevel:${TOOL_VERSION} kubectl-opslevel:latest 
cat << EOF > /usr/local/bin/kubectl-opslevel
#! /bin/sh
docker run -it --rm -v \$(pwd):/app -v ${HOME}/.kube:/.kube -e KUBECONFIG=/.kube/config --network=host kubectl-opslevel:latest \$@
EOF
chmod +x /usr/local/bin/kubectl-opslevel
```

<!---
TODO: Implement other methods


#### Homebrew


TODO: Need to Publish to Homebrew


```
brew update && brew install kubectl-opslevel
```

#### Windows


TODO: Chocolately?


1. Get `kubectl-opslevel-windows-amd64` from our [releases](https://github.com/opslevel/kubectl-opslevel/releases/latest).
2. Rename `kubectl-opslevel-windows-amd64` to `kubectl-opslevel.exe` and store it in a preferred path.
3. Make sure the location you choose is added to your Path environment variable.

-->

### Validate Install

Once you have the binary on your Path you can validate it works by running:

```sh
kubectl opslevel version
```

Example Output:

```sh
v0.4.0.0-g0d8107bdd043
```

The log format default is more human readable but if you want structured logs you can set the flag `--logFormat=JSON`

```json
{"level":"info","time":1620251466,"message":"Ensured tag 'imported = kubectl-opslevel' assigned to service: 'db'"}
```

## Documentation

You can read more about the tool and its configuration format in our [documentation](https://www.opslevel.com/docs/integrations/kubernetes/)
