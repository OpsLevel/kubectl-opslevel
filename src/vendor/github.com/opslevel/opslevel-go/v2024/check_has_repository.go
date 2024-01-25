package opslevel

func (client *Client) CreateCheckRepositoryIntegrated(input CheckRepositoryIntegratedCreateInput) (*Check, error) {
	var m struct {
		Payload CheckResponsePayload `graphql:"checkRepositoryIntegratedCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("CheckRepositoryIntegratedCreate"))
	return &m.Payload.Check, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) UpdateCheckRepositoryIntegrated(input CheckRepositoryIntegratedUpdateInput) (*Check, error) {
	var m struct {
		Payload CheckResponsePayload `graphql:"checkRepositoryIntegratedUpdate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("CheckRepositoryIntegratedUpdate"))
	return &m.Payload.Check, HandleErrors(err, m.Payload.Errors)
}
