package opslevel

type GitBranchProtectionCheckFragment struct{}

func (client *Client) CreateCheckGitBranchProtection(input CheckGitBranchProtectionCreateInput) (*Check, error) {
	var m struct {
		Payload CheckResponsePayload `graphql:"checkGitBranchProtectionCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("CheckGitBranchProtectionCreate"))
	return &m.Payload.Check, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) UpdateCheckGitBranchProtection(input CheckGitBranchProtectionUpdateInput) (*Check, error) {
	var m struct {
		Payload CheckResponsePayload `graphql:"checkGitBranchProtectionUpdate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("CheckGitBranchProtectionUpdate"))
	return &m.Payload.Check, HandleErrors(err, m.Payload.Errors)
}
