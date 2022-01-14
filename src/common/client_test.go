package common

import (
	"testing"

	"github.com/opslevel/opslevel-go"
	"github.com/rocktavious/autopilot"
)

func Test_ServiceNeedsUpdate_IsTrue_WhenInputDiffers(t *testing.T) {
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
		Name: "Test1",
	}
	input2 := opslevel.ServiceUpdateInput{
		Name:     "Test",
		Language: "Python",
		Tier:     "tier_2",
	}
	input3 := opslevel.ServiceUpdateInput{
		Name:     "Test",
		Language: "Python",
		Owner:    "platform",
	}
	// Act
	result1 := serviceNeedsUpdate(input1, &service)
	result2 := serviceNeedsUpdate(input2, &service)
	result3 := serviceNeedsUpdate(input3, &service)
	// Assert
	autopilot.Equals(t, true, result1)
	autopilot.Equals(t, true, result2)
	autopilot.Equals(t, true, result3)
}

func Test_ServiceNeedsUpdate_IsFalse_WhenInputMatches(t *testing.T) {
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
		Name:     "Test",
		Language: "Python",
		Tier:     "tier_1",
	}
	// Act
	result1 := serviceNeedsUpdate(input1, &service)
	result2 := serviceNeedsUpdate(input2, &service)
	result3 := serviceNeedsUpdate(input3, &service)
	// Assert
	autopilot.Equals(t, false, result1)
	autopilot.Equals(t, false, result2)
	autopilot.Equals(t, false, result3)
}
