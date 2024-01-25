package opslevel

import (
	"fmt"
	"slices"
)

type InfrastructureResourceSchema struct {
	Type   string `json:"type"`
	Schema JSON   `json:"schema" scalar:"true"`
}

type InfrastructureResourceSchemaConnection struct {
	Nodes      []InfrastructureResourceSchema
	PageInfo   PageInfo
	TotalCount int `graphql:"-"`
}

type InfrastructureResourceProviderData struct {
	AccountName  string `json:"accountName" graphql:"accountName"`
	ExternalURL  string `json:"externalUrl" graphql:"externalUrl"`
	ProviderName string `json:"providerName" graphql:"providerName"`
}

type InfrastructureResource struct {
	Id           string                             `json:"id"`
	Aliases      []string                           `json:"aliases"`
	Name         string                             `json:"name"`
	Schema       string                             `json:"type" graphql:"type @include(if: $all)"`
	ProviderType string                             `json:"providerResourceType" graphql:"providerResourceType @include(if: $all)"`
	ProviderData InfrastructureResourceProviderData `json:"providerData" graphql:"providerData @include(if: $all)"`
	Owner        EntityOwner                        `json:"owner" graphql:"owner @include(if: $all)"`
	OwnerLocked  bool                               `json:"ownerLocked" graphql:"ownerLocked @include(if: $all)"`
	ParsedData   JSON                               `json:"data" scalar:"true" graphql:"data @include(if: $all)"`
	Data         JSON                               `json:"rawData" scalar:"true" graphql:"rawData @include(if: $all)"`
}

type InfrastructureResourceConnection struct {
	Nodes      []InfrastructureResource
	PageInfo   PageInfo
	TotalCount int `graphql:"-"`
}

type InfraProviderInput struct {
	Account string `json:"account" yaml:"account" default:"Dev - 123456789"`
	Name    string `json:"name" yaml:"name" default:"Google"`
	Type    string `json:"type" yaml:"type" default:"BigQuery"`
	URL     string `json:"url" yaml:"url" default:"https://google.com"`
}

type InfraInput struct {
	Schema   string              `json:"schema" yaml:"schema" default:"Database"`
	Owner    *ID                 `json:"owner,omitempty" yaml:"owner,omitempty" default:"XXX_owner_id_XXX"`
	Provider *InfraProviderInput `json:"provider" yaml:"provider"`
	Data     *JSON               `json:"data" yaml:"data" default:"{\"name\":\"my-big-query\",\"engine\":\"BigQuery\",\"endpoint\":\"https://google.com\",\"replica\":false}"`
}

func (i *InfrastructureResource) GetTags(client *Client, variables *PayloadVariables) (*TagConnection, error) {
	var q struct {
		Account struct {
			InfrastructureResource struct {
				Tags TagConnection `graphql:"tags(after: $after, first: $first)"`
			} `graphql:"infrastructureResource(input: $infrastructureResource)"`
		}
	}
	if i.Id == "" {
		return nil, fmt.Errorf("Unable to get Tags, invalid InfrastructureResource id: '%s'", i.Id)
	}
	if variables == nil {
		variables = client.InitialPageVariablesPointer()
	}
	(*variables)["infrastructureResource"] = *NewIdentifier(string(i.Id))

	if err := client.Query(&q, *variables, WithName("InfrastructureResourceTags")); err != nil {
		return nil, err
	}
	for q.Account.InfrastructureResource.Tags.PageInfo.HasNextPage {
		(*variables)["after"] = q.Account.InfrastructureResource.Tags.PageInfo.End
		resp, err := i.GetTags(client, variables)
		if err != nil {
			return nil, err
		}
		// Add unique tags only
		for _, resp := range resp.Nodes {
			if !slices.Contains[[]Tag, Tag](q.Account.InfrastructureResource.Tags.Nodes, resp) {
				q.Account.InfrastructureResource.Tags.Nodes = append(q.Account.InfrastructureResource.Tags.Nodes, resp)
			}
		}
		q.Account.InfrastructureResource.Tags.PageInfo = resp.PageInfo
		q.Account.InfrastructureResource.Tags.TotalCount += resp.TotalCount
	}
	return &q.Account.InfrastructureResource.Tags, nil
}

func (i *InfrastructureResource) ResourceId() ID {
	return *NewID(i.Id)
}

func (i *InfrastructureResource) ResourceType() TaggableResource {
	return TaggableResourceInfrastructureresource
}

func (client *Client) CreateInfrastructure(input InfraInput) (*InfrastructureResource, error) {
	i := InfrastructureResourceInput{
		Schema: &InfrastructureResourceSchemaInput{Type: input.Schema},
		Data:   input.Data,
	}
	if input.Owner != nil {
		i.OwnerId = input.Owner
	}
	if input.Provider != nil {
		i.ProviderResourceType = &input.Provider.Type
		i.ProviderData = &InfrastructureResourceProviderDataInput{
			AccountName:  input.Provider.Account,
			ExternalUrl:  RefOf(input.Provider.URL),
			ProviderName: RefOf(input.Provider.Name),
		}
	}
	var m struct {
		Payload struct {
			InfrastructureResource InfrastructureResource
			Warnings               []OpsLevelWarnings
			Errors                 []OpsLevelErrors
		} `graphql:"infrastructureResourceCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": i,
		"all":   true,
	}
	err := client.Mutate(&m, v, WithName("InfrastructureResourceCreate"))
	// TODO: handle m.Payload.Warnings somehow
	return &m.Payload.InfrastructureResource, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) GetInfrastructure(identifier string) (*InfrastructureResource, error) {
	var q struct {
		Account struct {
			InfrastructureResource InfrastructureResource `graphql:"infrastructureResource(input: $input)"`
		}
	}
	v := PayloadVariables{
		"input": *NewIdentifier(identifier),
		"all":   true,
	}
	err := client.Query(&q, v, WithName("InfrastructureResourceGet"))
	if q.Account.InfrastructureResource.Id == "" {
		err = fmt.Errorf("InfrastructureResource with identifier '%s' not found", identifier)
	}
	return &q.Account.InfrastructureResource, HandleErrors(err, nil)
}

func (client *Client) ListInfrastructureSchemas(variables *PayloadVariables) (*InfrastructureResourceSchemaConnection, error) {
	var q struct {
		Account struct {
			InfrastructureResourceSchemas InfrastructureResourceSchemaConnection `graphql:"infrastructureResourceSchemas(after: $after, first: $first)"`
		}
	}
	if variables == nil {
		variables = client.InitialPageVariablesPointer()
	}
	if err := client.Query(&q, *variables, WithName("IntegrationList")); err != nil {
		return nil, err
	}
	for q.Account.InfrastructureResourceSchemas.PageInfo.HasNextPage {
		(*variables)["after"] = q.Account.InfrastructureResourceSchemas.PageInfo.End
		resp, err := client.ListInfrastructureSchemas(variables)
		if err != nil {
			return nil, err
		}
		q.Account.InfrastructureResourceSchemas.Nodes = append(q.Account.InfrastructureResourceSchemas.Nodes, resp.Nodes...)
		q.Account.InfrastructureResourceSchemas.PageInfo = resp.PageInfo
	}
	q.Account.InfrastructureResourceSchemas.TotalCount = len(q.Account.InfrastructureResourceSchemas.Nodes)
	return &q.Account.InfrastructureResourceSchemas, nil
}

func (client *Client) ListInfrastructure(variables *PayloadVariables) (*InfrastructureResourceConnection, error) {
	var q struct {
		Account struct {
			InfrastructureResource InfrastructureResourceConnection `graphql:"infrastructureResources(after: $after, first: $first)"`
		}
	}
	if variables == nil {
		variables = client.InitialPageVariablesPointer()
		(*variables)["all"] = true
	}
	if err := client.Query(&q, *variables, WithName("IntegrationList")); err != nil {
		return nil, err
	}
	for q.Account.InfrastructureResource.PageInfo.HasNextPage {
		(*variables)["after"] = q.Account.InfrastructureResource.PageInfo.End
		resp, err := client.ListInfrastructure(variables)
		if err != nil {
			return nil, err
		}
		q.Account.InfrastructureResource.Nodes = append(q.Account.InfrastructureResource.Nodes, resp.Nodes...)
		q.Account.InfrastructureResource.PageInfo = resp.PageInfo
	}
	q.Account.InfrastructureResource.TotalCount = len(q.Account.InfrastructureResource.Nodes)
	return &q.Account.InfrastructureResource, nil
}

func (client *Client) UpdateInfrastructure(identifier string, input InfraInput) (*InfrastructureResource, error) {
	i := InfrastructureResourceInput{
		Data: input.Data,
	}
	if input.Owner != nil {
		i.OwnerId = input.Owner
	}
	if input.Provider != nil {
		i.ProviderResourceType = &input.Provider.Type
		i.ProviderData = &InfrastructureResourceProviderDataInput{
			AccountName:  input.Provider.Account,
			ExternalUrl:  RefOf(input.Provider.URL),
			ProviderName: RefOf(input.Provider.Name),
		}
	}
	var m struct {
		Payload struct {
			InfrastructureResource InfrastructureResource
			Warnings               []OpsLevelWarnings
			Errors                 []OpsLevelErrors
		} `graphql:"infrastructureResourceUpdate(infrastructureResource: $identifier, input: $input)"`
	}
	v := PayloadVariables{
		"identifier": *NewIdentifier(identifier),
		"input":      i,
		"all":        true,
	}
	err := client.Mutate(&m, v, WithName("InfrastructureResourceUpdate"))
	// TODO: handle m.Payload.Warnings somehow
	return &m.Payload.InfrastructureResource, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) DeleteInfrastructure(identifier string) error {
	var m struct {
		Payload struct {
			Errors []OpsLevelErrors `graphql:"errors"`
		} `graphql:"infrastructureResourceDelete(resource: $input)"`
	}
	v := PayloadVariables{
		"input": *NewIdentifier(identifier),
	}
	err := client.Mutate(&m, v, WithName("InfrastructureResourceDelete"))
	return HandleErrors(err, m.Payload.Errors)
}
