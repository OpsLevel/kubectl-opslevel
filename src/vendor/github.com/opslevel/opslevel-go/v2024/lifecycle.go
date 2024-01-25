package opslevel

type Lifecycle struct {
	Alias       string
	Description string
	Id          ID
	Index       int
	Name        string
}

func (client *Client) ListLifecycles() ([]Lifecycle, error) {
	var q struct {
		Account struct {
			Lifecycles []Lifecycle
		}
	}
	err := client.Query(&q, nil, WithName("LifecycleList"))
	return q.Account.Lifecycles, HandleErrors(err, nil)
}
