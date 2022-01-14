package common

import (
	"testing"

	"github.com/opslevel/opslevel-go"
	"github.com/rocktavious/autopilot"
)

func TestServiceNeedsUpdate(t *testing.T) {
	// Arrange
	service := opslevel.Service{
		ServiceId: opslevel.ServiceId{
			Id: opslevel.NewID("XXX"),
		},
		Name:        "Test",
		Description: "Hello World",
		Language:    "Python",
		Tier: opslevel.Tier{
			Alias: "tier_1",
		},
	}
	input1 := opslevel.ServiceUpdateInput{
		Id: "XXX",
	}
	input2 := opslevel.ServiceUpdateInput{
		Name: "Test",
	}
	input3 := opslevel.ServiceUpdateInput{
		Name: "Test1",
	}
	input4 := opslevel.ServiceUpdateInput{
		Name:     "Test",
		Language: "Python",
		Tier:     "tier_1",
	}
	input5 := opslevel.ServiceUpdateInput{
		Name:     "Test",
		Language: "Python",
		Tier:     "tier_2",
	}
	input6 := opslevel.ServiceUpdateInput{
		Name:     "Test",
		Language: "Python",
		Owner:    "platform",
	}
	// Act
	result1 := serviceNeedsUpdate(input1, &service)
	result2 := serviceNeedsUpdate(input2, &service)
	result3 := serviceNeedsUpdate(input3, &service)
	result4 := serviceNeedsUpdate(input4, &service)
	result5 := serviceNeedsUpdate(input5, &service)
	result6 := serviceNeedsUpdate(input6, &service)
	// Assert
	autopilot.Equals(t, false, result1)
	autopilot.Equals(t, false, result2)
	autopilot.Equals(t, true, result3)
	autopilot.Equals(t, false, result4)
	autopilot.Equals(t, true, result5)
	autopilot.Equals(t, true, result6)
}
