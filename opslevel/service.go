package opslevel

import (
	"github.com/shurcooL/graphql"
)

type Service struct {
	Aliases []graphql.String
	//CheckStats
	//Dependencies
	//Dependents
	Description graphql.String
	Framework graphql.String
	Id graphql.ID
	Language graphql.String
	Lifecycle Lifecycle
	Name graphql.String
	Owner Team
	Product graphql.String
	//Repositories
	//Tags
	Tier Tier
	//Tools
}

type ServiceCreateInput struct {
	Name graphql.String `json:"name"`
	Product graphql.String `json:"product,omitempty"`
	Description graphql.String `json:"description,omitempty"`
	Languague graphql.String `json:"language,omitempty"`
	Framework graphql.String `json:"framework,omitempty"`
	Tier graphql.String `json:"tierAlias,omitempty"`
	Owner graphql.String `json:"ownerAlias,omitempty"`
	Lifecycle graphql.String `json:"lifecycleAlias,omitempty"`
}

type ServiceUpdateInput struct {
	Id graphql.ID `json:"id,omitempty"`
	Alias graphql.String `json:"alias,omitempty"`
	Name graphql.String `json:"name,omitempty"`
	Product graphql.String `json:"product,omitempty"`
	Descripition graphql.String `json:"description,omitempty"`
	Languague graphql.String `json:"languague,omitempty"`
	Framework graphql.String `json:"framework,omitempty"`
	Tier graphql.String `json:"tierAlias,omitempty"`
	Owner graphql.String `json:"ownerAlias,omitempty"`
	Lifecycle graphql.String `json:"lifecycleAlias,omitempty"`}

type ServiceDeleteInput struct {
	Id graphql.ID `json:"id,omitempty"`
	Alias graphql.String `json:"alias,omitempty"`
}

func (client *Client) GetServiceWithAlias(alias string) (*Service, error) {
	var q struct {
		Account struct {
			Service Service `graphql:"service(alias: $service)"`
		}
	}
	v := PayloadVariables{
		"service": graphql.String(alias),
	}
	if err := client.Query(&q, v); err != nil {
		return nil, err
	}
	return &q.Account.Service, nil
}

func (client *Client) GetServiceWithId(id string) (*Service, error) {
	var q struct {
		Account struct {
			Service Service `graphql:"service(id: $service)"`
		}
	}
	v := PayloadVariables{
		"service": graphql.ID(id),
	}
	if err := client.Query(&q, v); err != nil {
		return nil, err
	}
	return &q.Account.Service, nil
}

func (client *Client) CreateService(input ServiceCreateInput) (*Service, error) {
	var m struct {
		Payload struct {
			Service Service
			Errors []OpsLevelErrors
		} `graphql:"serviceCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	if err := client.Mutate(&m, v); err != nil {
		return nil, err
	}
	return &m.Payload.Service, FormatErrors(m.Payload.Errors)
}

func (client *Client) UpdateService(input ServiceUpdateInput) (*Service, error) {
	var m struct {
		Payload struct {
			Service Service
			Errors []OpsLevelErrors
		} `graphql:"serviceUpdate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	if err := client.Mutate(&m, v); err != nil {
		return nil, err
	}
	return &m.Payload.Service, FormatErrors(m.Payload.Errors)
}

func (client *Client) DeleteService(input ServiceDeleteInput) error {
	var m struct {
		Payload struct {
			Id graphql.ID `graphql:"deletedServiceId"`
			Alias graphql.String `graphql:"deletedServiceAlias"`
			Errors []OpsLevelErrors `graphql:"errors"`
		} `graphql:"serviceDelete(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	if err := client.Mutate(&m, v); err != nil {
		return err
	}
	return FormatErrors(m.Payload.Errors)
}