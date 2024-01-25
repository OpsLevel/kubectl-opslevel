package opslevel

func (client *Client) CreateCheckServiceConfiguration(input CheckServiceConfigurationCreateInput) (*Check, error) {
	var m struct {
		Payload CheckResponsePayload `graphql:"checkServiceConfigurationCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("CheckServiceConfigurationCreate"))
	return &m.Payload.Check, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) UpdateCheckServiceConfiguration(input CheckServiceConfigurationUpdateInput) (*Check, error) {
	var m struct {
		Payload CheckResponsePayload `graphql:"checkServiceConfigurationUpdate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("CheckServiceConfigurationUpdate"))
	return &m.Payload.Check, HandleErrors(err, m.Payload.Errors)
}
