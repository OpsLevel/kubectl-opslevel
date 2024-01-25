package opslevel

type Tier struct {
	Alias       string
	Description string
	Id          ID
	Index       int
	Name        string
}

func (client *Client) ListTiers() ([]Tier, error) {
	var q struct {
		Account struct {
			Tiers []Tier
		}
	}
	err := client.Query(&q, nil, WithName("TierList"))
	return q.Account.Tiers, HandleErrors(err, nil)
}
