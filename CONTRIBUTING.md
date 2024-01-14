# Contributing

1. [Local Development](#local-development)

## Local Development

You can set up a safe and convenient environment for testing by using [minikube](https://minikube.sigs.k8s.io/docs/start/).
`brew install minikube`

We use [taskfiles as our build tool](https://taskfile.dev/) to install dependencies. `brew install go-task`

### Compile

```
task setup && cd src && go build . && ./kubectl-opslevel
```

> [!TIP]
> Your local version can be called only by using `./kubectl-opslevel` in only that directory. If you want to use it everywhere, you can copy it to /usr/local/bin and delete your currently installed version. You can verify which one is installed by using `kubectl plugin list`.

### Start minikube

```
minikube start
```

> [!CAUTION]
> Beyond this step, make sure your `kubectl` is set to your minikube cluster AND that your `OPSLEVEL_API_TOKEN` is for a sandbox environment.

### Optional: create k8s resources for testing

You can skip both of these steps if you want to use your own test data.

```
$ cat deploy.yaml | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-web-staging
  labels:
    app: nginx
  annotations:
    opslevel.com/description: nginx web deployed for staging.
    opslevel.com/owner: platform
    opslevel.com/lifecycle: beta
    opslevel.com/tier: tier_4
    opslevel.com/tags.cloud: aws
    opslevel.com/tags.has_pii: "false"
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        environment: staging
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
```

```
$ cat opslevel-k8s.yaml
version: 1.2.0
service:
    import:
        - selector:
            apiVersion: apps/v1
            kind: Deployment
            excludes:
                - .metadata.namespace == "kube-system"
                - .metadata.annotations."opslevel.com/ignore"
          opslevel:
            name: .metadata.name
            description: .metadata.annotations."opslevel.com/description"
            owner: .metadata.annotations."opslevel.com/owner"
            aliases:
                - '"\(.metadata.namespace)-\(.metadata.name)"'
            tags:
                assign:
                    - '{"k8s_last_scanned": now | strftime("%Y-%m-%dT%H:%M:%SZ") }'
                    - '{"k8s_created": .metadata.creationTimestamp}'
                    - '.metadata.annotations | to_entries |  map(select(.key | startswith("opslevel.com/tags"))) | map({(.key | split(".")[2]): .value})'
                    - .metadata.labels
    collect:
        - selector:
            apiVersion: apps/v1
            kind: Deployment
            excludes:
                - .metadata.namespace == "kube-system"
                - .metadata.annotations."opslevel.com/ignore"
```

### Cleanup

`minikube stop`
