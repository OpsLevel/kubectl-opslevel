package opslevel

import (
	"github.com/shurcooL/graphql"
)

type Tier struct {
	Alias       graphql.String
	Description graphql.String
	Id          graphql.String
	Index       graphql.Int
	Name        graphql.String
}

func (client *Client) ListTiers() ([]Tier, error) {
	var q struct {
		Account struct {
			Tiers []Tier
		}
	}
	if err := client.Query(&q, nil); err != nil {
		return []Tier{}, err
	}
	return q.Account.Tiers, nil
}
