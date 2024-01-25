package opslevel

import (
	"fmt"
)

type ScorecardId struct {
	Aliases []string `graphql:"aliases"`
	Id      ID       `graphql:"id"`
}

type Scorecard struct {
	ScorecardId

	AffectsOverallServiceLevels bool        `graphql:"affectsOverallServiceLevels"`
	Description                 string      `graphql:"description"` // optional
	Filter                      Filter      `graphql:"filter"`      // optional
	Name                        string      `graphql:"name"`
	Owner                       EntityOwner `graphql:"owner"`
	PassingChecks               int         `graphql:"passingChecks"`
	ServiceCount                int         `graphql:"serviceCount"`
	ChecksCount                 int         `graphql:"totalChecks"`
}

type ScorecardConnection struct {
	Nodes      []Scorecard `graphql:"nodes"`
	PageInfo   PageInfo    `graphql:"pageInfo"`
	TotalCount int         `graphql:"totalCount"`
}

func (client *Client) CreateScorecard(input ScorecardInput) (*Scorecard, error) {
	var m struct {
		Payload struct {
			Scorecard Scorecard        `graphql:"scorecard"`
			Errors    []OpsLevelErrors `graphql:"errors"`
		} `graphql:"scorecardCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("ScorecardCreate"))
	return &m.Payload.Scorecard, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) GetScorecard(input string) (*Scorecard, error) {
	var q struct {
		Account struct {
			Scorecard Scorecard `graphql:"scorecard(input: $input)"`
		}
	}
	v := PayloadVariables{
		"input": *NewIdentifier(input),
	}
	err := client.Query(&q, v, WithName("ScorecardGet"))
	if q.Account.Scorecard.Id == "" {
		err = fmt.Errorf("Scorecard with ID or Alias matching '%s' not found", input)
	}
	return &q.Account.Scorecard, HandleErrors(err, nil)
}

func (client *Client) ListScorecards(variables *PayloadVariables) (*ScorecardConnection, error) {
	var q struct {
		Account struct {
			Scorecards ScorecardConnection `graphql:"scorecards(after: $after, first: $first)"`
		}
	}
	if variables == nil {
		variables = client.InitialPageVariablesPointer()
	}
	if err := client.Query(&q, *variables, WithName("ScorecardsList")); err != nil {
		return nil, err
	}
	for q.Account.Scorecards.PageInfo.HasNextPage {
		(*variables)["after"] = q.Account.Scorecards.PageInfo.End
		resp, err := client.ListScorecards(variables)
		if err != nil {
			return nil, err
		}
		q.Account.Scorecards.Nodes = append(q.Account.Scorecards.Nodes, resp.Nodes...)
		q.Account.Scorecards.PageInfo = resp.PageInfo
		q.Account.Scorecards.TotalCount = len(q.Account.Scorecards.Nodes)
	}
	return &q.Account.Scorecards, nil
}

func (client *Client) UpdateScorecard(identifier string, input ScorecardInput) (*Scorecard, error) {
	var m struct {
		Payload struct {
			Scorecard Scorecard        `graphql:"scorecard"`
			Errors    []OpsLevelErrors `graphql:"errors"`
		} `graphql:"scorecardUpdate(scorecard: $scorecard, input: $input)"`
	}
	scorecard := *NewIdentifier(identifier)
	v := PayloadVariables{
		"scorecard": scorecard,
		"input":     input,
	}
	err := client.Mutate(&m, v, WithName("ScorecardUpdate"))
	return &m.Payload.Scorecard, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) DeleteScorecard(identifier string) (*ID, error) {
	var m struct {
		Payload struct {
			DeletedScorecardId ID               `graphql:"deletedScorecardId"`
			Errors             []OpsLevelErrors `graphql:"errors"`
		} `graphql:"scorecardDelete(input: $input)"`
	}
	input := *NewIdentifier(identifier)
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("ScorecardDelete"))
	return &m.Payload.DeletedScorecardId, HandleErrors(err, m.Payload.Errors)
}
