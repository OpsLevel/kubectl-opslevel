package opslevel

import (
	"github.com/shurcooL/graphql"
)

type User struct {
	Name  graphql.String
	Email graphql.String
}

type Contact struct {
	DisplayName graphql.String
	Address     graphql.String
}

type ContactInput struct {
	Type        string `json:"type,omitEmpty"`
	DisplayName string `json:"displayName,omitEmpty"`
	Address     string `json:"address,omitEmpty"`
}

type Team struct {
	Id               graphql.String
	Name             graphql.String
	Responsibilities graphql.String
	Manager          User
	Contacts         []Contact
}

type TeamCreateInput struct {
	Name             string         `json:"name"`
	ManagerEmail     string         `json:"managerEmail,omitempty"`
	Responsibilities string         `json:"responsibilities,omitempty"`
	Contacts         []ContactInput `json:"contacts,omitempty"`
}

type TeamUpdateInput struct {
	Id               graphql.ID `json:"id,omitempty"`
	Alias            string     `json:"alias,omitempty"`
	Name             string     `json:"name,omitempty"`
	ManagerEmail     string     `json:"managerEmail,omitempty"`
	Responsibilities string     `json:"responsibilities,omitempty"`
}

type TeamDeleteInput struct {
	Id    graphql.ID `json:"id,omitempty"`
	Alias string     `json:"alias,omitempty"`
}

//#region Create

func (client *Client) CreateTeam(input TeamCreateInput) (*Team, error) {
	var m struct {
		Payload struct {
			Team   Team
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

//#endregion

//#region Retrieve

func (client *Client) GetTeamWithAlias(alias string) (*Team, error) {
	var q struct {
		Account struct {
			Team Team `graphql:"team(alias: $team)"`
		}
	}
	v := PayloadVariables{
		"team": alias,
	}
	if err := client.Query(&q, v); err != nil {
		return nil, err
	}
	return &q.Account.Team, nil
}

func (client *Client) GetTeamWithId(id graphql.ID) (*Team, error) {
	var q struct {
		Account struct {
			Team Team `graphql:"team(id: $team)"`
		}
	}
	v := PayloadVariables{
		"team": id,
	}
	if err := client.Query(&q, v); err != nil {
		return nil, err
	}
	return &q.Account.Team, nil
}

func (client *Client) GetTeamCount() (int, error) {
	var q struct {
		Account struct {
			Teams struct {
				TotalCount graphql.Int
			}
		}
	}
	if err := client.Query(&q, nil); err != nil {
		return 0, err
	}
	return int(q.Account.Teams.TotalCount), nil
}

type ListTeamsQuery struct {
	Account struct {
		Teams struct {
			Nodes    []Team
			PageInfo PageInfo
		} `graphql:"teams(after: $after, first: $first)"`
	}
}

func (q *ListTeamsQuery) Query(client *Client) error {
	var subQ ListTeamsQuery
	v := PayloadVariables{
		"after": q.Account.Teams.PageInfo.End,
		"first": graphql.Int(100),
	}
	if err := client.Query(&subQ, v); err != nil {
		return err
	}
	if subQ.Account.Teams.PageInfo.HasNextPage {
		subQ.Query(client)
	}
	for _, team := range subQ.Account.Teams.Nodes {
		q.Account.Teams.Nodes = append(q.Account.Teams.Nodes, team)
	}
	return nil
}

func (client *Client) ListTeams() ([]Team, error) {
	q := ListTeamsQuery{}
	if err := q.Query(client); err != nil {
		return []Team{}, err
	}
	return q.Account.Teams.Nodes, nil
}

//#endregion

//#region Update

func (client *Client) UpdateTeam(input TeamUpdateInput) (*Team, error) {
	var m struct {
		Payload struct {
			Team   Team
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

//#endregion

//#region Delete

func (client *Client) DeleteTeamWithAlias(alias string) error {
	var m struct {
		Payload struct {
			Id     graphql.ID       `graphql:"deletedTeamId"`
			Alias  graphql.String   `graphql:"deletedTeamAlias"`
			Errors []OpsLevelErrors `graphql:"errors"`
		} `graphql:"teamDelete(input: $input)"`
	}
	v := PayloadVariables{
		"input": TeamDeleteInput{
			Alias: alias,
		},
	}
	if err := client.Mutate(&m, v); err != nil {
		return err
	}
	return FormatErrors(m.Payload.Errors)
}

func (client *Client) DeleteTeamWithId(id graphql.ID) error {
	var m struct {
		Payload struct {
			Id     graphql.ID       `graphql:"deletedTeamId"`
			Alias  graphql.String   `graphql:"deletedTeamAlias"`
			Errors []OpsLevelErrors `graphql:"errors"`
		} `graphql:"teamDelete(input: $input)"`
	}
	v := PayloadVariables{
		"input": TeamDeleteInput{
			Id: id,
		},
	}
	if err := client.Mutate(&m, v); err != nil {
		return err
	}
	return FormatErrors(m.Payload.Errors)
}

//#endregion
