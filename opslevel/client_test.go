package opslevel

import (
	"flag"
	"os"
	"testing"

	"github.com/rocktavious/autopilot"

	"github.com/shurcooL/graphql"
)

func TestMain(m *testing.M) {
	flag.Parse()
	teardown := autopilot.Setup()
	defer teardown()
	os.Exit(m.Run())
}

func TestClientQuery(t *testing.T) {
	// Arrange
	url := autopilot.RegisterEndpoint("/query_account", "query_account.json")
	client := NewClient("X", SetURL(url))
	var q struct {
		Account struct {
			Id graphql.ID
		}
	}
	var v map[string]interface{}
	// Act
	err := client.Query(&q, v)
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, "1234", q.Account.Id)
}
