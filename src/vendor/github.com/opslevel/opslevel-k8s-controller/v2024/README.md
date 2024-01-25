<p align="center">
    <a href="https://github.com/OpsLevel/opslevel-k8s-controller/blob/main/LICENSE" alt="License">
        <img src="https://img.shields.io/github/license/OpsLevel/opslevel-k8s-controller.svg" /></a>
    <a href="http://golang.org" alt="Made With Go">
        <img src="https://img.shields.io/github/go-mod/go-version/OpsLevel/opslevel-k8s-controller" /></a>
    <a href="https://GitHub.com/OpsLevel/opslevel-k8s-controller/releases/" alt="Release">
        <img src="https://img.shields.io/github/v/release/OpsLevel/opslevel-k8s-controller?include_prereleases" /></a>
    <a href="https://GitHub.com/OpsLevel/opslevel-k8s-controller/issues/" alt="Issues">
        <img src="https://img.shields.io/github/issues/OpsLevel/opslevel-k8s-controller.svg" /></a>
    <a href="https://github.com/OpsLevel/opslevel-k8s-controller/graphs/contributors" alt="Contributors">
        <img src="https://img.shields.io/github/contributors/OpsLevel/opslevel-k8s-controller" /></a>
    <a href="https://github.com/OpsLevel/opslevel-k8s-controller/pulse" alt="Activity">
        <img src="https://img.shields.io/github/commit-activity/m/OpsLevel/opslevel-k8s-controller" /></a>
	<a href="https://codecov.io/gh/OpsLevel/opslevel-k8s-controller">
  		<img src="https://codecov.io/gh/OpsLevel/opslevel-k8s-controller/branch/main/graph/badge.svg?token=GHQHRIJ9UW"/></a>
    <a href="https://dependabot.com/" alt="Dependabot">
        <img src="https://badgen.net/badge/Dependabot/enabled/green?icon=dependabot" /></a>
    <a href="https://pkg.go.dev/github.com/opslevel/opslevel-k8s-controller/v2024" alt="Go Reference">
        <img src="https://pkg.go.dev/badge/github.com/opslevel/opslevel.svg" /></a>
</p>


[![Overall](https://img.shields.io/endpoint?style=flat&url=https%3A%2F%2Fapp.opslevel.com%2Fapi%2Fservice_level%2FDEmrX2bjoPualtC4Pri1YvjydDrza6V1V5srMvcZNbQ)](https://app.opslevel.com/services/opslevel-k8s-controller/maturity-report)

# opslevel-k8s-controller
A utility library for easily making and running k8s controllers

# Installation

```bash
go get github.com/opslevel/opslevel-k8s-controller/v2024
```

Then to create a k8s controller you can simply do

```go
selector := opslevel_k8s_controller.K8SSelector{
    ApiVersion: "apps/v1",
    Kind: "Deployment",
    Excludes: []string{`.metadata.namespace == "kube-system"`}
}
resync := time.Hour*24
batch := 500
runOnce := false
controller, err := opslevel_k8s_controller.NewK8SController(selector, resync, batch, runOnce)
if err != nil {
    //... Handle error ...
}
callback := func(items []interface{}) {
    for _, item := range items {
        // ... Process K8S Resource ...
    }
}
controller.OnAdd = callback
controller.OnUpdate = callback
controller.Start()
```

Because of the way the selector works you can easily target any k8s resource in your cluster and you have the power of JQ
to exclude resources that might match the expression.