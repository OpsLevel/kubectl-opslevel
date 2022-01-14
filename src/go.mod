module github.com/opslevel/kubectl-opslevel

go 1.16

require (
	github.com/alecthomas/jsonschema v0.0.0-20210526225647-edb03dcab7bc
	github.com/creasty/defaults v1.5.1
	github.com/go-logr/logr v0.4.0
	github.com/google/go-cmp v0.5.5
	github.com/imdario/mergo v0.3.10 // indirect
	github.com/opslevel/opslevel-go v0.4.2
	github.com/rocktavious/autopilot v0.1.5 // indirect
	github.com/rs/zerolog v1.26.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	go.uber.org/automaxprocs v1.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	k8s.io/klog/v2 v2.9.0
)

// Uncomment for local development
// replace github.com/opslevel/opslevel-go => ../../opslevel-go

replace github.com/opslevel/kubectl-opslevel/common => ./common

replace github.com/opslevel/kubectl-opslevel/cmd => ./cmd

replace github.com/opslevel/kubectl-opslevel/config => ./config

replace github.com/opslevel/kubectl-opslevel/jq => ./jq

replace github.com/opslevel/kubectl-opslevel/k8utils => ./k8utils
