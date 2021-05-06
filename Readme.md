# kubectl-opslevel

<!---
TOC Generated with https://github.com/ekalinin/github-markdown-toc
-->
## Table of Contents

* [Prerequisite](#prerequisite)
* [Installation](#installation)
      * [MacOS](#macos)
      * [Linux](#linux)
      * [Docker](#docker)
  * [Validate Install](#validate-install)
* [Working With The Configuration File](#working-with-the-configuration-file)
* [Preview](#preview)
* [Import](#import)

## Prerequisite

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [jq](https://stedolan.github.io/jq/download/)
- [OpsLevel API Token](https://app.opslevel.com/api_tokens)

## Installation

#### MacOS

```
curl -Lo kubectl-opslevel.tar.gz https://github.com/opslevel/kubectl-opslevel/releases/download/$(curl -s https://api.github.com/repos/opslevel/kubectl-opslevel/releases/latest | grep tag_name | cut -d '"' -f 4)/kubectl-opslevel-darwin-amd64.tar.gz
tar -xzvf kubectl-opslevel.tar.gz  
rm kubectl-opslevel.tar.gz
chmod +x kubectl-opslevel
sudo mv kubectl-opslevel /usr/local/bin/kubectl-opslevel
```

#### Linux

```
curl -Lo kubectl-opslevel https://github.com/opslevel/kubectl-opslevel/releases/download/$(curl -s https://api.github.com/repos/opslevel/kubectl-opslevel/releases/latest | grep tag_name | cut -d '"' -f 4)/kubectl-opslevel-linux-amd64.tar.gz
tar -xzvf kubectl-opslevel.tar.gz  
rm kubectl-opslevel.tar.gz
chmod +x kubectl-opslevel
sudo mv kubectl-opslevel /usr/local/bin/kubectl-opslevel
```

#### Docker

```
docker pull public.ecr.aws/opslevel/kubectl-opslevel:latest
```

Then run the following script to inject a shim into your `/usr/local/bin` so you can use the binary like its downloaded natively - it will just be running in a docker container.

```
cat << EOF > /usr/local/bin/kubectl-opslevel
#! /bin/sh
docker run -it --rm -w /mounted -v \$(pwd):/mounted public.ecr.aws/opslevel/kubectl-opslevel:latest \$@
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

```
kubectl opslevel version
```

Example Output:

```
{"level":"info","time":1620251466,"message":"v0.1.1-0-gc52681db6b33"}
```

## Working With The Configuration File

The tool is driven by a configuration file that allows you to map data from kubernetes resource into OpsLevel fields.  Here is a simple example that maps a deployment's metadata name to an OpsLevel service name:

<!---
TODO: Would be great to read this from a static file in the repo/wiki?
-->

```yaml
service:
  import:
  - selector:
      kind: deployment
    opslevel:
      name: .metadata.name
```

You can also generate a sample config to act as a starting point with the following command:

```
kubectl opslevel config sample > ./opslevel-k8s.yaml
```

The sample configuration file provides a bunch of sane defaults to help get you started but you likely will need to tweak it to work for your organization specific needs.  The following are some advanced examples of things you might want to do.

Filter out unwanted keys from your labels to tags mapping (in this case excluding labels that start with "flux"):

```
service:
  import:
  - selector:
      kind: deployment
    opslevel:
      name: '"\(.metadata.name)-\(.metadata.namespace)"'
      tags:
      - .metadata.labels | to_entries | map(select(.key |startswith("flux") | not)) | from_entries
```

Target `Ingress` Resources where each host rule is attached as a tool to cataloge the domains used:

```
service:
  import:
  - selector:
      kind: ingress
    opslevel:
      name: '"\(.metadata.name)-\(.metadata.namespace)"'
      tags:
      - {"imported": "kubectl-opslevel"}
      tools:
      - '.spec.rules | map({"cateogry":"other","displayName":.host,"url": .host})'
```

## Preview

The primary iteration loop of the tool resides in tweaking the configuration file and running the `service preview` command to view data that represents what the tool will do (think of this as a dry-run or terraform plan)

```
kubectl opslevel service preview -c ./opslevel-k8s.yaml
```

Once you are happy with the full output you can move onto the actual import process.

## Import

Once you are ready to import data into your Opslevel account run the following:
*(insert your OpsLevel API Token)*:

```
 OL_APITOKEN=XXXX kubectl opslevel service import -c ./opslevel-k8s.yaml
```

This command may take a few minutes to run so please be patient while it works.  In the meantime you can open a browser to your [OpsLevel account](https://app.opslevel.com/) and view the newly generated services.