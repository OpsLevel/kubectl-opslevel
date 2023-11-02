package common_test

import (
	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/rocktavious/autopilot/v2023"
	"testing"
)

func TestGetConfig(t *testing.T) {
	simple, err := common.GetConfig(common.ConfigSimple)
	autopilot.Ok(t, err)
	sample, err := common.GetConfig(common.ConfigSample)
	autopilot.Ok(t, err)

	autopilot.Equals(t, ".metadata.namespace", simple.Service.Import[0].OpslevelConfig.Owner)
	autopilot.Equals(t, ".metadata.annotations.\"opslevel.com/owner\"", sample.Service.Import[0].OpslevelConfig.Owner)
}
