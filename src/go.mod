module github.com/opslevel/kubectl-opslevel

go 1.16

require (
	github.com/alecthomas/jsonschema v0.0.0-20210526225647-edb03dcab7bc
	github.com/creasty/defaults v1.5.1
	github.com/go-logr/logr v0.3.0
	github.com/google/go-cmp v0.5.5
	github.com/googleapis/gnostic v0.5.1 // indirect
	github.com/imdario/mergo v0.3.10 // indirect
	github.com/opslevel/opslevel-go v0.4.0
	github.com/rs/zerolog v1.25.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	go.uber.org/automaxprocs v1.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
	k8s.io/klog/v2 v2.4.0
	k8s.io/utils v0.0.0-20210111153108-fddb29f9d009 // indirect
)

// Uncomment for local development
// replace github.com/opslevel/opslevel-go => ../../opslevel-go
