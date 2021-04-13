package opslevel

import (
	"github.com/shurcooL/graphql"
)

type User struct {
	Name graphql.String
	Email graphql.String
}

type Contact struct {
	DisplayName graphql.String
	Address graphql.String
}

type Team struct {
	Id graphql.String
	Name graphql.String
	Responsibilities graphql.String
	Manager User
	Contacts []Contact
}

type TeamCreateInput struct {
	Name graphql.String `json:"name"`
	ManagerEmail graphql.String `json:"managerEmail,omitempty"`
	Responsibilities graphql.String `json:"responsibilities,omitempty"`
	Contacts []ContactInput `json:"contacts,omitempty"`
}

type ContactInput struct {
	Type graphql.String `json:"type,omitEmpty"`
    DisplayName graphql.String `json:"displayName,omitEmpty"`
	Address graphql.String `json:"address,omitEmpty"`
}

type TeamUpdateInput struct {
	Id graphql.ID `json:"id,omitempty"`
	Alias graphql.String `json:"alias,omitempty"`
	Name graphql.String `json:"name,omitempty"`
	ManagerEmail graphql.String `json:"managerEmail,omitempty"`
	Responsibilities graphql.String `json:"responsibilities,omitempty"`
}

type TeamDeleteInput struct {
	Id graphql.ID `json:"id,omitempty"`
	Alias graphql.String `json:"alias,omitempty"`
}

func (client *Client) GetTeamWithAlias(alias string) (*Team, error) {
	var q struct {
		Account struct {
			Team Team `graphql:"team(alias: $team)"`
		}
	}
	v := PayloadVariables{
		"team": graphql.String(alias),
	}
	if err := client.Query(&q, v); err != nil {
		return nil, err
	}
	return &q.Account.Team, nil
}

func (client *Client) GetTeamWithId(id string) (*Team, error) {
	var q struct {
		Account struct {
			Team Team `graphql:"team(id: $team)"`
		}
	}
	v := PayloadVariables{
		"team": graphql.ID(id),
	}
	if err := client.Query(&q, v); err != nil {
		return nil, err
	}
	return &q.Account.Team, nil
}

func (client *Client) CreateTeam(input TeamCreateInput) (*Team, error) {
	var m struct {
		Payload struct {
			Team Team
			Errors []OpsLevelErrors
		} `graphql:"teamCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	if err := client.Mutate(&m, v); err != nil {
		return nil, err
	}
	return &m.Payload.Team, FormatErrors(m.Payload.Errors)
}

func (client *Client) UpdateTeam(input TeamUpdateInput) (*Team, error) {
	var m struct {
		Payload struct {
			Team Team
			Errors []OpsLevelErrors
		} `graphql:"teamUpdate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	if err := client.Mutate(&m, v); err != nil {
		return nil, err
	}
	return &m.Payload.Team, FormatErrors(m.Payload.Errors)
}

func (client *Client) DeleteTeam(input TeamDeleteInput) error {
	var m struct {
		Payload struct {
			Id graphql.ID `graphql:"deletedTeamId"`
			Alias graphql.String `graphql:"deletedTeamAlias"`
			Errors []OpsLevelErrors `graphql:"errors"`
		} `graphql:"teamDelete(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	if err := client.Mutate(&m, v); err != nil {
		return err
	}
	return FormatErrors(m.Payload.Errors)
}