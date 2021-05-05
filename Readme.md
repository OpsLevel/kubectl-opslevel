# kubectl-opslevel

[[_TOC_]]

## Prerequisite

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [jq](https://stedolan.github.io/jq/download/)
- [OpsLevel API Token](https://app.opslevel.com/api_tokens)

## Install

#### Homebrew

<!---
TODO: Need to Publish to Homebrew
-->

```
brew update && brew install kubectl-opslevel
```

#### MacOS

<!---
TODO: Need to build and publish non tar'd binaries
-->

```
curl -Lo kubectl-opslevel https://github.com/opslevel/kubectl-opslevel/releases/download/$(curl -s https://api.github.com/repos/opslevel/kubectl-opslevel/releases/latest | grep tag_name | cut -d '"' -f 4)/kubectl-opslevel-darwin-amd64
chmod +x kubectl-opslevel
sudo mv kubectl-opslevel /usr/local/bin/kubectl-opslevel
```

#### Linux

```
curl -Lo kubectl-opslevel https://github.com/opslevel/kubectl-opslevel/releases/download/$(curl -s https://api.github.com/repos/opslevel/kubectl-opslevel/releases/latest | grep tag_name | cut -d '"' -f 4)/kubectl-opslevel-linux-amd64
chmod +x kubectl-opslevel
sudo mv kubectl-opslevel /usr/local/bin/kubectl-opslevel
```

#### Windows

<!---
TODO: Chocolately?
-->

1. Get `kubectl-opslevel-windows-amd64` from our [releases](https://github.com/opslevel/kubectl-opslevel/releases/latest).
2. Rename `kubectl-opslevel-windows-amd64` to `kubectl-opslevel.exe` and store it in a preferred path.
3. Make sure the location you choose is added to your Path environment variable.

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

This command may take a bit to run so be patient (we will be further imporving the UX of this output in the future) 

Once the command is complete open up your [OpsLevel account](https://app.opslevel.com/) and view all the newly generated services