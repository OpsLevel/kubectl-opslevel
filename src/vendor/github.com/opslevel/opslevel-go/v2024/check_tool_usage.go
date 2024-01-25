package opslevel

type ToolUsageCheckFragment struct {
	ToolCategory         ToolCategory `graphql:"toolCategory"`
	ToolNamePredicate    *Predicate   `graphql:"toolNamePredicate"`
	ToolUrlPredicate     *Predicate   `graphql:"toolUrlPredicate"`
	EnvironmentPredicate *Predicate   `graphql:"environmentPredicate"`
}

func (client *Client) CreateCheckToolUsage(input CheckToolUsageCreateInput) (*Check, error) {
	var m struct {
		Payload CheckResponsePayload `graphql:"checkToolUsageCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("CheckToolUsageCreate"))
	return &m.Payload.Check, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) UpdateCheckToolUsage(input CheckToolUsageUpdateInput) (*Check, error) {
	var m struct {
		Payload CheckResponsePayload `graphql:"checkToolUsageUpdate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("CheckToolUsageUpdate"))
	return &m.Payload.Check, HandleErrors(err, m.Payload.Errors)
}
