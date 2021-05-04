package opslevel

import (
	"github.com/shurcooL/graphql"
)

type AliasCreateInput struct {
	Alias   string     `json:"alias"`
	OwnerId graphql.ID `json:"ownerId"`
}

//#region Create

func (client *Client) CreateAliases(ownerId graphql.ID, aliases []string) []string {
	var output []string
	for _, alias := range aliases {
		input := AliasCreateInput{
			Alias:   alias,
			OwnerId: ownerId,
		}
		result, err := client.CreateAlias(input)
		if err != nil {
			// TODO: log warning about failed create?
		}
		for _, resultAlias := range result {
			output = append(output, string(resultAlias))
		}
	}
	// TODO: need to treat this like a HashSet to deduplicate
	return output
}

func (client *Client) CreateAlias(input AliasCreateInput) ([]string, error) {
	var m struct {
		Payload struct {
			Aliases []graphql.String
			OwnerId graphql.String
			Errors  []OpsLevelErrors
		} `graphql:"aliasCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	if err := client.Mutate(&m, v); err != nil {
		return nil, err
	}
	output := make([]string, len(m.Payload.Aliases))
	for i, item := range m.Payload.Aliases {
		output[i] = string(item)
	}
	return output, FormatErrors(m.Payload.Errors)
}

//#endregion
