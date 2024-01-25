package opslevel

import (
	"fmt"
	"strings"
)

// #region Create
// TODO: make sure duplicate aliases throw an error that we can catch
func (client *Client) CreateAliases(ownerId ID, aliases []string) ([]string, error) {
	var output []string
	var errors []string
	for _, alias := range aliases {
		input := AliasCreateInput{
			Alias:   alias,
			OwnerId: ownerId,
		}
		result, err := client.CreateAlias(input)
		if err != nil {
			errors = append(errors, err.Error())
		}
		for _, resultAlias := range result {
			output = append(output, string(resultAlias))
		}
	}
	output = removeDuplicates(output)
	if len(errors) > 0 {
		return output, fmt.Errorf(strings.Join(errors, "\n"))
	} else {
		return output, nil
	}
}

func (client *Client) CreateAlias(input AliasCreateInput) ([]string, error) {
	var m struct {
		Payload struct {
			Aliases []string
			OwnerId string
			Errors  []OpsLevelErrors
		} `graphql:"aliasCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("AliasCreate"))
	output := make([]string, len(m.Payload.Aliases))
	for i, item := range m.Payload.Aliases {
		output[i] = string(item)
	}
	return output, HandleErrors(err, m.Payload.Errors)
}

//#endregion

//#region Delete

func (client *Client) DeleteInfraAlias(alias string) error {
	return client.DeleteAlias(AliasDeleteInput{
		Alias:     alias,
		OwnerType: AliasOwnerTypeEnumInfrastructureResource,
	})
}

func (client *Client) DeleteServiceAlias(alias string) error {
	return client.DeleteAlias(AliasDeleteInput{
		Alias:     alias,
		OwnerType: AliasOwnerTypeEnumService,
	})
}

func (client *Client) DeleteTeamAlias(alias string) error {
	return client.DeleteAlias(AliasDeleteInput{
		Alias:     alias,
		OwnerType: AliasOwnerTypeEnumTeam,
	})
}

func (client *Client) DeleteAlias(input AliasDeleteInput) error {
	var m struct {
		Payload struct {
			Alias  string `graphql:"deletedAlias"`
			Errors []OpsLevelErrors
		} `graphql:"aliasDelete(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("AliasDelete"))
	return HandleErrors(err, m.Payload.Errors)
}

//#endregion
