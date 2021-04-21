package opslevel

import (
	"github.com/shurcooL/graphql"
)

type AliasCreateInput struct {
	Alias graphql.String `json:"alias"`
	OwnerId graphql.String `json:"ownerId"`
}

//#region Create

func (client *Client) CreateAliases(ownerId string, aliases []string) []graphql.String {
	var output []graphql.String
	for _, alias := range aliases {
		input := AliasCreateInput{
			Alias: graphql.String(alias),
			OwnerId: graphql.String(ownerId),
		}
		result, err := client.CreateAlias(input)
		if err != nil {
			// TODO: log warning about failed create?
		}
		for _, resultAlias := range result {
			output = append(output, resultAlias)
		}
	}
	// TODO: need to treat this like a HashSet to deduplicate
	return output
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
	return m.Payload.Aliases, FormatErrors(m.Payload.Errors)
}

//#endregion
