package opslevel

import (
	"github.com/shurcooL/graphql"
)

type AliasCreateInput struct {
	Alias graphql.String `json:"alias"`
	OwnerId graphql.String `json:"ownerId"`
}

//#region Create

func (client *Client) CreateAliases(ownerId string, aliases []string) {
	for _, alias := range aliases {
		input := AliasCreateInput{
			Alias: alias,
			OwnerId: ownerId,
		}
		if err := client.CreateAlias(input); err != nil {
			// TODO: log warning about failed create?
		}
		// TODO: should we append all aliases and return a final result []graphql.String?
	}
}

func (client *Client) CreateAlias(input AliasCreateInput) ([]graphql.String, error) {
	var m struct {
		Payload struct {
			Aliases []graphql.String
			OwnerId graphql.String
			Errors []OpsLevelErrors
		} `graphql:"aliasCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	if err := client.Mutate(&m, v); err != nil {
		return nil, err
	}
	return &m.Payload.Aliases, FormatErrors(m.Payload.Errors)
}

//#endregion
