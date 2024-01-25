package opslevel

type ServicePropertyCheckFragment struct {
	Property  ServicePropertyTypeEnum `graphql:"serviceProperty"`
	Predicate *Predicate              `graphql:"propertyValuePredicate"`
}

func (client *Client) CreateCheckServiceProperty(input CheckServicePropertyCreateInput) (*Check, error) {
	var m struct {
		Payload CheckResponsePayload `graphql:"checkServicePropertyCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("CheckServicePropertyCreate"))
	return &m.Payload.Check, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) UpdateCheckServiceProperty(input CheckServicePropertyUpdateInput) (*Check, error) {
	var m struct {
		Payload CheckResponsePayload `graphql:"checkServicePropertyUpdate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("CheckServicePropertyUpdate"))
	return &m.Payload.Check, HandleErrors(err, m.Payload.Errors)
}
