package common_test

import (
	"testing"

	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/rocktavious/autopilot/v2023"
)

func TestParseConfig(t *testing.T) {
	simple, err := common.ParseConfig(common.ConfigSimple)
	autopilot.Ok(t, err)
	sample, err := common.ParseConfig(common.ConfigSample)
	autopilot.Ok(t, err)

	autopilot.Equals(t, ".metadata.namespace", simple.Service.Import[0].OpslevelConfig.Owner)
	autopilot.Equals(t, ".metadata.annotations.\"opslevel.com/owner\"", sample.Service.Import[0].OpslevelConfig.Owner)
}
