# kubectl-opslevel
---

## Quickstart

Download the latest binary from the [release page](https://github.com/OpsLevel/kubectl-opslevel/releases) untar the binary and place it in `/usr/local/bin/`

Then generate a configuration file - you can get a starting sample using the following command

```
kubectl opslevel config sample > ./opslevel-k8s.yaml
```

You and preview the data the tool will use to generate services in your OpsLevel account with:

```
kubectl opslevel service preview -c ./opslevel-k8s.yaml
```

Tweak the `./opslevel-k8s.yaml` and continue re-running the `service preview` command until you are happy with the full output.

Once you are ready to import the data from the preview command into your OpsLevel account you run the following (modify the following command with your own OpsLevel account API Token)

```
 OL_APITOKEN=XXXX kubectl opslevel service import -c ./opslevel-k8s.yaml
```

This command may take a bit to run so be patient (we will be further imporving the UX of this output in the future) 

Once the command is complete open up your [OpsLevel account](https://app.opslevel.com/) and view all the newly generated services