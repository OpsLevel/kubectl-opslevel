<a name="unreleased"></a>
## [Unreleased]


<a name="v0.2.0"></a>
## [v0.2.0] - 2021-05-22
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

[Unreleased]: https://github.com/OpsLevel/kubectl-opslevel/compare/v0.2.0...HEAD
[v0.2.0]: https://github.com/OpsLevel/kubectl-opslevel/compare/v0.1.3...v0.2.0
[v0.1.3]: https://github.com/OpsLevel/kubectl-opslevel/compare/v0.1.2...v0.1.3
[v0.1.2]: https://github.com/OpsLevel/kubectl-opslevel/compare/v0.1.1...v0.1.2
[v0.1.1]: https://github.com/OpsLevel/kubectl-opslevel/compare/v0.1.0...v0.1.1
[v0.1.0]: https://github.com/OpsLevel/kubectl-opslevel/compare/v0.0.1...v0.1.0
