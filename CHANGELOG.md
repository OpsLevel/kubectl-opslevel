<a name="unreleased"></a>
## [Unreleased]


<a name="0.4.4"></a>
## [0.4.4] - 2021-07-22
### Bugfix
- fix index lookup error when a service field is not configured in a selector resulting in a empty lookup array
- fix logic around aggregating services so that multiple import selectors work


<a name="v0.4.3"></a>
## [v0.4.3] - 2021-07-21
### Bugfix
- properly handle jq expressions that return ‘empty’ rather then null now that we batch expression operations

### Feature
- upgrade jq from 1.4 to 1.6 ([#58](https://github.com/OpsLevel/kubectl-opslevel/issues/58))


<a name="v0.4.2"></a>
## [v0.4.2] - 2021-07-14
### Docs
- simplify install instructions to use homebrew

### Feature
- switch exclude filters to use JQ multiParser to increase performance on filtering resources
- use new multiParse JQ processing when processing the service fields to increase performance when working with 100’s of resources in Kubernetes
- implement ability to run JQ statements on arrays of resources to speed up processing when retrieving 100’s of resources from Kubernetes


<a name="v0.4.1"></a>
## [v0.4.1] - 2021-07-10
### Refactor
- switch to goreleaser for publishing

### Reverts
- Revert "Extract Documentation to OpsLevel Website"
- Extract Documentation to OpsLevel Website


<a name="v0.4.0"></a>
## [v0.4.0] - 2021-06-25
### Docs
- update readme to use new opslevel AWS ECR alias
- update sample configuration file in readme to reflect config version 1.1.0
- switch demo to asciinema
- move readme documentation to OpsLevel website

### Feature
- add ability to filter resources based on jq expressions that return truthy
- add config validation messages for new selector excludes format
- add config selector validation to help user upgrade from 1.0.0 -> 1.1.0 with helpful messages
- support targeting ANY k8s resource in the cluster

### Refactor
- bump config version to support query any resource kind in the cluster


<a name="v0.3.1-beta.1"></a>
## [v0.3.1-beta.1] - 2021-06-13
### Bugfix
- protect against empty alias value
- protect against empty key or value in tag assign and create lists
- remove overlapped keys that exist in parsed `tags.assign` if key name also exists in `tags.create`
- adjust golang binary inside docker container to have same version as when built outside the container

### Docs
- add install instructions for using docker container

### Feature
- add ability to reconcile linked service repository display name
- add ability to assign repositories to services
- add kubernetes deploy instructions that link to the helm chart


<a name="v0.3.0.beta1"></a>
## [v0.3.0.beta1] - 2021-05-29
### Bugfix
- support all auth methods for kubeconfig by importing the auth package
- load kubeconfig using client-go ConfigLoadingRules to support multi file KUBECONFIG settings

### Docs
- remove 'include_prereleases' from the release badge
- fix badge link to go.mod
- add a reminder to install the tool first in the quickstart section
- rework the readme to have a better flow for new users

### Feature
- support tags that are updated and tags that are append only
- Provide a way to exclude a namespace(s) via the selector
- add publishing of docker image to AWS public OpsLevel ECR
- support printing out comments as part of the config sample command to help explain aspects about the configuration file
- support a version field in the configuration file to allow for making backwards incompatible changes
- improve the `service import` output to include an indentation character for better readability
- improve the output that surrounds the `service preview` data to better instruct the user what to do next

### Refactor
- switch final docker image to use ubuntu instead of golang
- k8sutils resource type handler functions use []byte to reduce code duplication
- cleanup folder structure to improve docker build caching
- change version command to just print out the version rather then using a log statement


<a name="v0.2.0"></a>
## [v0.2.0] - 2021-05-22
### Bugfix
- flag `—api-token` did not work but the environment variable did

### Chore
- upgrade to opslevel-go 0.2.0
- add badges to readme
- set MIT license

### Docs
- added "beta" status callout in readme
- add terminalizer demo file and gif
- add quick description to beginning of readme after badges

### Feature
- add ability to list aliases for tool categories `account tools`
- add a more simple sample configuration using the command with an additional flag `config sample —simple`
- enrich output during client validation when 401 unauthorized

### Refactor
- move opslevel client creation to a common function to always validate for an API key
- move HasTool check to opslevel-go and validate existance with environment not url to allow for updating the url of a tool


<a name="v0.1.3"></a>
## [v0.1.3] - 2021-05-16
### Chore
- update opslevel-go to 0.1.3

### Docs
- add quickstart to readme
- remove docker install instructions from the readme

### Feature
- validate jq is installed and on the path - if not log a message
- add a log statement when finished importing data for a given service
- skip over aliases and tags that already exist on the service
- Simple reconcile of existing service entries
- set a more opinionated aliases configuration with a default

### Refactor
- bump the log statement for skipping a tool down to debug level
- find services via the aliases list rather then using name
- improve error output when unable to cache tiers, lifecycle, teams
- change default log format to be more human readable
- Improve service import log output to explain what steps were success or failed ([#6](https://github.com/OpsLevel/kubectl-opslevel/issues/6))
- extract opslevel api client library to opslevel-go ([#8](https://github.com/OpsLevel/kubectl-opslevel/issues/8))


<a name="v0.1.2"></a>
## [v0.1.2] - 2021-05-08
### Chore
- Init Changelog

### Docs
- Update readme TOC
- Write a more robust readme to help with getting started

### Feature
- Support attaching Lifecycle, Tier and Owner on newly created service entries
- Implement listing Team aliases available in your account
- Implement listing Tier aliases available in your account
- Implement listing Lifecycle aliases available in your account
- Implement Account command for viewing data about your OpsLevel account
- Support creating tags on newly created service entries
- Support creating tools on newly created service entries


<a name="v0.1.1"></a>
## [v0.1.1] - 2021-04-22

<a name="v0.1.0"></a>
## [v0.1.0] - 2021-04-22

<a name="v0.0.1"></a>
## v0.0.1 - 2021-03-25

[Unreleased]: https://github.com/OpsLevel/kubectl-opslevel/compare/0.4.4...HEAD
[0.4.4]: https://github.com/OpsLevel/kubectl-opslevel/compare/v0.4.3...0.4.4
[v0.4.3]: https://github.com/OpsLevel/kubectl-opslevel/compare/v0.4.2...v0.4.3
[v0.4.2]: https://github.com/OpsLevel/kubectl-opslevel/compare/v0.4.1...v0.4.2
[v0.4.1]: https://github.com/OpsLevel/kubectl-opslevel/compare/v0.4.0...v0.4.1
[v0.4.0]: https://github.com/OpsLevel/kubectl-opslevel/compare/v0.3.1-beta.1...v0.4.0
[v0.3.1-beta.1]: https://github.com/OpsLevel/kubectl-opslevel/compare/v0.3.0.beta1...v0.3.1-beta.1
[v0.3.0.beta1]: https://github.com/OpsLevel/kubectl-opslevel/compare/v0.2.0...v0.3.0.beta1
[v0.2.0]: https://github.com/OpsLevel/kubectl-opslevel/compare/v0.1.3...v0.2.0
[v0.1.3]: https://github.com/OpsLevel/kubectl-opslevel/compare/v0.1.2...v0.1.3
[v0.1.2]: https://github.com/OpsLevel/kubectl-opslevel/compare/v0.1.1...v0.1.2
[v0.1.1]: https://github.com/OpsLevel/kubectl-opslevel/compare/v0.1.0...v0.1.1
[v0.1.0]: https://github.com/OpsLevel/kubectl-opslevel/compare/v0.0.1...v0.1.0
