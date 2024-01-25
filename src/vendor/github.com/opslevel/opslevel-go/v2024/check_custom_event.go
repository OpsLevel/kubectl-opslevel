package opslevel

type CustomEventCheckFragment struct {
	Integration      IntegrationId `graphql:"integration"`
	PassPending      bool          `graphql:"passPending"`
	ResultMessage    string        `graphql:"resultMessage"`
	ServiceSelector  string        `graphql:"serviceSelector"`
	SuccessCondition string        `graphql:"successCondition"`
}

func (client *Client) CreateCheckCustomEvent(input CheckCustomEventCreateInput) (*Check, error) {
	var m struct {
		Payload CheckResponsePayload `graphql:"checkCustomEventCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("CheckCustomEventCreate"))
	return &m.Payload.Check, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) UpdateCheckCustomEvent(input CheckCustomEventUpdateInput) (*Check, error) {
	var m struct {
		Payload CheckResponsePayload `graphql:"checkCustomEventUpdate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("CheckCustomEventUpdate"))
	return &m.Payload.Check, HandleErrors(err, m.Payload.Errors)
}
