module github.com/opslevel/kubectl-opslevel

go 1.18

require (
	github.com/alecthomas/jsonschema v0.0.0-20210526225647-edb03dcab7bc
	github.com/creasty/defaults v1.6.0
	github.com/go-logr/logr v1.2.3
	github.com/go-resty/resty/v2 v2.7.0
	github.com/google/go-cmp v0.5.8
	github.com/opslevel/opslevel-go/v2022 v2022.10.22
	github.com/rocktavious/autopilot v0.1.5
	github.com/rs/zerolog v1.28.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	go.uber.org/automaxprocs v1.5.1
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/apimachinery v0.25.4
	k8s.io/client-go v0.25.3
	k8s.io/klog/v2 v2.80.1
)

require (
	cloud.google.com/go v0.97.0 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.27 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.20 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful/v3 v3.8.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.2.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/gosimple/slug v1.13.1 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/iancoleman/orderedmap v0.0.0-20190318233801-ac98e3ecb4b0 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pelletier/go-toml v1.9.3 // indirect
	github.com/relvacode/iso8601 v1.1.0 // indirect
	github.com/shurcooL/graphql v0.0.0-20220606043923-3cf50f8a0a29 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	golang.org/x/crypto v0.0.0-20220315160706-3147a52a75dd // indirect
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b // indirect
	golang.org/x/oauth2 v0.0.0-20220608161450-d0670ef3b1eb // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20220210224613-90d013bbcef8 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.25.3 // indirect
	k8s.io/kube-openapi v0.0.0-20220803162953-67bda5d908f1 // indirect
	k8s.io/utils v0.0.0-20220728103510-ee6ede2d64ed // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

// Uncomment for local development
// replace github.com/opslevel/opslevel-go => ./submodules/opslevel-go
