module github.com/opslevel/kubectl-opslevel

go 1.16

require (
	github.com/creasty/defaults v1.5.1
	github.com/go-logr/logr v0.3.0
	github.com/google/go-cmp v0.5.5
	github.com/googleapis/gnostic v0.5.1 // indirect
	github.com/imdario/mergo v0.3.10 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/opslevel/opslevel-go v0.2.0
	github.com/rs/zerolog v1.21.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
	k8s.io/klog/v2 v2.4.0
	k8s.io/utils v0.0.0-20210111153108-fddb29f9d009 // indirect
)

// Uncomment for local development
// replace github.com/opslevel/opslevel-go => ../opslevel-go
