package opslevel

import (
	"fmt"

	"github.com/gosimple/slug"
)

type Predicate struct {
	Type  PredicateTypeEnum `graphql:"type"`
	Value string            `graphql:"value"`
}

type FilterId struct {
	Id   ID
	Name string
}

type Filter struct {
	Connective ConnectiveEnum
	HtmlURL    string
	FilterId
	Predicates []FilterPredicate
}

type FilterPredicate struct {
	Key           PredicateKeyEnum  `json:"key" yaml:"key" default:"repository_ids"`
	KeyData       string            `json:"keyData,omitempty" yaml:"keyData,omitempty" default:"null"`
	Type          PredicateTypeEnum `json:"type" yaml:"type" default:"equals"`
	Value         string            `json:"value,omitempty" yaml:"value,omitempty" default:"1"`
	CaseSensitive *bool             `json:"caseSensitive,omitempty" yaml:"caseSensitive,omitempty" default:"false"`
}

type FilterConnection struct {
	Nodes      []Filter
	PageInfo   PageInfo
	TotalCount int
}

func (self *Filter) Alias() string {
	return slug.Make(self.Name)
}

//#region Create

func (client *Client) CreateFilter(input FilterCreateInput) (*Filter, error) {
	var m struct {
		Payload struct {
			Filter Filter
			Errors []OpsLevelErrors
		} `graphql:"filterCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("FilterCreate"))
	return &m.Payload.Filter, HandleErrors(err, m.Payload.Errors)
}

//#endregion

//#region Retrieve

func (client *Client) GetFilter(id ID) (*Filter, error) {
	var q struct {
		Account struct {
			Filter Filter `graphql:"filter(id: $id)"`
		}
	}
	v := PayloadVariables{
		"id": id,
	}
	err := client.Query(&q, v, WithName("FilterGet"))
	if q.Account.Filter.Id == "" {
		err = fmt.Errorf("Filter with ID '%s' not found!", id)
	}
	return &q.Account.Filter, HandleErrors(err, nil)
}

func (client *Client) ListFilters(variables *PayloadVariables) (*FilterConnection, error) {
	var q struct {
		Account struct {
			Filters FilterConnection `graphql:"filters(after: $after, first: $first)"`
		}
	}
	if variables == nil {
		variables = client.InitialPageVariablesPointer()
	}
	if err := client.Query(&q, *variables, WithName("FilterList")); err != nil {
		return nil, err
	}
	for q.Account.Filters.PageInfo.HasNextPage {
		(*variables)["after"] = q.Account.Filters.PageInfo.End
		resp, err := client.ListFilters(variables)
		if err != nil {
			return nil, err
		}
		q.Account.Filters.Nodes = append(q.Account.Filters.Nodes, resp.Nodes...)
		q.Account.Filters.PageInfo = resp.PageInfo
		q.Account.Filters.TotalCount += resp.TotalCount
	}
	return &q.Account.Filters, nil
}

//#endregion

//#region Update

func (client *Client) UpdateFilter(input FilterUpdateInput) (*Filter, error) {
	var m struct {
		Payload struct {
			Filter Filter
			Errors []OpsLevelErrors
		} `graphql:"filterUpdate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("FilterUpdate"))
	return &m.Payload.Filter, HandleErrors(err, m.Payload.Errors)
}

//#endregion

//#region Delete

func (client *Client) DeleteFilter(id ID) error {
	var m struct {
		Payload struct {
			Id     ID `graphql:"deletedId"`
			Errors []OpsLevelErrors
		} `graphql:"filterDelete(input: $input)"`
	}
	v := PayloadVariables{
		"input": DeleteInput{Id: id},
	}
	err := client.Mutate(&m, v, WithName("FilterDelete"))
	return HandleErrors(err, m.Payload.Errors)
}

//#endregion
