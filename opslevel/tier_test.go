package opslevel

import (
	"testing"

	"github.com/rocktavious/autopilot"
)

func TestListTiers(t *testing.T) {
	// Arrange
	url := autopilot.RegisterEndpoint("/query_tiers", "query_tiers.json")
	client := NewClient("X", SetURL(url))
	// Act
	result, err := client.ListTiers()
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, 4, len(result))
	autopilot.Equals(t, []byte("tier_1"), []byte(result[0].Alias))
}
