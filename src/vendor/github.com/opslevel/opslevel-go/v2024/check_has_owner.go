package opslevel

type ServiceOwnershipCheckFragment struct {
	RequireContactMethod *bool        `graphql:"requireContactMethod"`
	ContactMethod        *ContactType `graphql:"contactMethod"`
	TeamTagKey           string       `graphql:"tagKey"`
	TeamTagPredicate     *Predicate   `graphql:"tagPredicate"`
}

func (client *Client) CreateCheckServiceOwnership(input CheckServiceOwnershipCreateInput) (*Check, error) {
	var m struct {
		Payload CheckResponsePayload `graphql:"checkServiceOwnershipCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("CheckServiceOwnershipCreate"))
	return &m.Payload.Check, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) UpdateCheckServiceOwnership(input CheckServiceOwnershipUpdateInput) (*Check, error) {
	var m struct {
		Payload CheckResponsePayload `graphql:"checkServiceOwnershipUpdate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("CheckServiceOwnershipUpdate"))
	return &m.Payload.Check, HandleErrors(err, m.Payload.Errors)
}
