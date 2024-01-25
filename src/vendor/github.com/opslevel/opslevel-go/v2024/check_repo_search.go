package opslevel

type RepositorySearchCheckFragment struct {
	FileExtensions        []string  `graphql:"fileExtensions"`
	FileContentsPredicate Predicate `graphql:"fileContentsPredicate"`
}

func (client *Client) CreateCheckRepositorySearch(input CheckRepositorySearchCreateInput) (*Check, error) {
	var m struct {
		Payload CheckResponsePayload `graphql:"checkRepositorySearchCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("CheckRepositorySearchCreate"))
	return &m.Payload.Check, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) UpdateCheckRepositorySearch(input CheckRepositorySearchUpdateInput) (*Check, error) {
	var m struct {
		Payload CheckResponsePayload `graphql:"checkRepositorySearchUpdate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("CheckRepositorySearchUpdate"))
	return &m.Payload.Check, HandleErrors(err, m.Payload.Errors)
}
