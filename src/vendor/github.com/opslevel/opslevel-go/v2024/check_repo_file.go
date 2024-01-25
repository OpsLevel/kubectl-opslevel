package opslevel

type RepositoryFileCheckFragment struct {
	DirectorySearch       bool       `graphql:"directorySearch"`
	Filepaths             []string   `graphql:"filePaths"`
	FileContentsPredicate *Predicate `graphql:"fileContentsPredicate"`
	UseAbsoluteRoot       bool       `graphql:"useAbsoluteRoot"`
}

func (client *Client) CreateCheckRepositoryFile(input CheckRepositoryFileCreateInput) (*Check, error) {
	var m struct {
		Payload CheckResponsePayload `graphql:"checkRepositoryFileCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("CheckRepositoryFileCreate"))
	return &m.Payload.Check, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) UpdateCheckRepositoryFile(input CheckRepositoryFileUpdateInput) (*Check, error) {
	var m struct {
		Payload CheckResponsePayload `graphql:"checkRepositoryFileUpdate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("CheckRepositoryFileUpdate"))
	return &m.Payload.Check, HandleErrors(err, m.Payload.Errors)
}
