package config_test

import (
	"github.com/opslevel/kubectl-opslevel/pkg/config"
	"github.com/rocktavious/autopilot"
	"testing"
)

func TestConfig(t *testing.T) {
	// Arrange
	simple := config.GetSample(true)
	sample := config.GetSample(false)
	custom := []byte(`version: "1.3.0"
service:
  import:
    - selector: &deployments
        apiVersion: apps/v1 # only supports resources found in 'kubectl api-resources --verbs="get,list"'
        kind: Deployment
      opslevel:
        name: .metadata.name
  collect:
    - selector:
        <<: *deployments
`)
	// Act
	cfg, err := config.NewConfig(custom)
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, simple, string(config.Simple))
	autopilot.Equals(t, sample, string(config.Sample))
	autopilot.Equals(t, config.ConfigCurrentVersion, cfg.Version)
	autopilot.Equals(t, "Deployment", cfg.Service.Import[0].SelectorConfig.Kind)
	autopilot.Equals(t, ".metadata.name", cfg.Service.Import[0].OpslevelConfig.Name)
	autopilot.Equals(t, "Deployment", cfg.Service.Collect[0].SelectorConfig.Kind)
}

func TestConfigErrors(t *testing.T) {
	// Arrange
	not_yaml := []byte(`{"version": "1.3.0"
service:
`)
	// Act
	cfg, err := config.NewConfig(not_yaml)
	// Assert
	autopilot.Assert(t, err != nil, "Parsed Invalid YAML")
	autopilot.Assert(t, cfg == nil, "Returned non-nil Config")
}
