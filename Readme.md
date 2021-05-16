Table of Contents
=================

<!--ts-->
   * [Prerequisite](#prerequisite)
   * [Installation](#installation)
         * [MacOS](#macos)
         * [Linux](#linux)
      * [Validate Install](#validate-install)
   * [Quickstart](#quickstart)
   * [Working With The Configuration File](#working-with-the-configuration-file)
      * [Sample Configuration Explained](#sample-configuration-explained)
         * [Lifecycle &amp; Tier &amp; Owner](#lifecycle--tier--owner)
         * [Aliases](#service-aliases)
         * [Tags](#tags)
         * [Tools](#tools)
   * [Preview](#preview)
   * [Import](#import)
<!--te-->

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

<!---
TODO: Implement other methods

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
4:35PM INF v0.1.1-0-gc52681db6b33
```

The log format default is more human readable but if you want structured logs you can set the flag `--logFormat=JSON`

```
{"level":"info","time":1620251466,"message":"v0.1.1-0-gc52681db6b33"}
```

## Quickstart

```
cat << EOF > ./opslevel-k8s.yaml
service:
  import:
  - selector:
      kind: deployment
    opslevel:
      # Fields can be a hardcoded value or jq expressions that generate strings
      name: .metadata.name
      product: .metadata.namespace
      language: "python"
      framework: "django"
      aliases:
      - '"k8s:\(.metadata.name)-\(.metadata.namespace)"'
      # For more complex data types you can use jq to build json objects or grab them right out of a k8s yaml field
      tags:
      - '{"imported": "kubectl-opslevel"}'
      - .spec.template.metadata.labels
EOF
kubectl opslevel service preview
 OL_APITOKEN=XXXX kubectl opslevel service import
```

You can also run this tool inside your kubernetes cluster as a job to reconcile the data with your OpsLevel account perodically.  Please read our [aliases](#aliases) section as we use these to properly find and reconcile the data so it is important you choose something unique.

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

### Sample Configuration Explained

In the sample configuration there are number of sane default jq expression set that should help you get started quickly.  Here we will breakdown some of the more advanced expressions to further your understanding and hopefully give you an idea of the power jq brings to data collection.

#### Lifecycle & Tier & Owner

The `Lifecycle`, `Tier` and `Owner` fields in the configuration are only validated upon `service import` not during `service preview`.  Valid values for these fields need to match the `Alias` for these resources in your OpsLevel account.  To view the valid aliases for these resources you can run the commands `account lifecycles`, `account tiers` and `account teams`.

#### Service Aliases

Aliases on services are very import to get right so that we can effectively look up the service later for reconciliation. Because of that we have only included 1 example which should be fairly unique to most kubernetes setup.  

*This example is also generated by default if no aliases are configured so that we can support lookup for reconcilation when you provide no other aliases.*

```
      aliases:
      - 'k8s:"\(.metadata.name)-\(.metadata.namespace)"'
```

The example concatenates togeather the resource name and namespace with a prefix of `k8s:` to create a unique alias.

#### Tags

In the tags section there are 4 example expressions that show you different ways to build the key/value payload for attaching tag entries to your service

```
      tags:
      - '{"imported": "kubectl-opslevel"}'
      - '.metadata.annotations | to_entries |  map(select(.key | startswith("opslevel.com/tags"))) | map({(.key | split(".")[2]): .value})'
      - .metadata.labels
      - .spec.template.metadata.labels
```

The first example shows how to hardcode a tag entry.  In this case we are denoting that the imported service came from this tool.

The second example leverages a convention to capture `1..N` tags.  The jq expression is looking for kubernetes annotations using the following format `opslevel.com/tags.<key>: <value>` and here is an example:

```
  annotations:
    opslevel.com/tags.hello: world
```

The third and fourth examples extract the `labels` applied to the kubernetes resource directly into your OpsLevel service's tags

#### Tools

In the tools section there are 2 example expressions that show you how to build the necessary payload for attaching tools entries to your service.

```
      tools:
      - '{"category": "other", "displayName": "my-cool-tool", "url": .metadata.annotations."example.com/my-cool-tool"}
        | if .url then . else empty end'
      - '.metadata.annotations | to_entries |  map(select(.key | startswith("opslevel.com/tools")))
        | map({"category": .key | split(".")[2], "displayName": .key | split(".")[3],
        "url": .value})'
```

The first example shows you the 3 required fields - `category` , `displayName` and `url` - but the expression hardcodes the values for `category` and `displayName` leaving you to specify where the `url` field comes from.

The second example leverages a convention to capture `1..N` tools.  The jq expression is looking for kubernetes annotations using the following format `opslevel.com/tools.<category>.<displayName>: <url>` and here is a example:

```
  annotations:
    opslevel.com/tools.logs.datadog: https://app.datadoghq.com/logs
```

## Preview

The primary iteration loop of the tool resides in tweaking the configuration file and running the `service preview` command to view data that represents what the tool will do (think of this as a dry-run or terraform plan)

```
kubectl opslevel service preview -c ./opslevel-k8s.yaml
```

Once you are happy with the full output you can move onto the actual import process.

*NOTE: this step does not validate any of the data with Opslevel - fields that are references to other things (IE: Tier, Lifecycle, Owner, etc) are not validated at this point and might cause a warning message during import* 

## Import

Once you are ready to import data into your Opslevel account run the following:
*(insert your OpsLevel API Token)*:

```
 OL_APITOKEN=XXXX kubectl opslevel service import -c ./opslevel-k8s.yaml
```

This command may take a few minutes to run so please be patient while it works.  In the meantime you can open a browser to your [OpsLevel account](https://app.opslevel.com/) and view the newly generated services.
