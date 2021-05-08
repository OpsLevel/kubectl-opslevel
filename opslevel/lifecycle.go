package opslevel

import (
	"github.com/shurcooL/graphql"
)

type Lifecycle struct {
	Alias       graphql.String
	Description graphql.String
	Id          graphql.ID
	Index       graphql.Int
	Name        graphql.String
}

func (client *Client) ListLifecycles() ([]Lifecycle, error) {
	var q struct {
		Account struct {
			Lifecycles []Lifecycle
		}
	}
	if err := client.Query(&q, nil); err != nil {
		return []Lifecycle{}, err
	}
	return q.Account.Lifecycles, nil
}
