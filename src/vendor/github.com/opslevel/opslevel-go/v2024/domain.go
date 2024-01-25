package opslevel

import (
	"fmt"
	"slices"
)

type DomainId Identifier

type Domain struct {
	DomainId
	Name        string      `graphql:"name"`
	Description string      `graphql:"description"`
	HTMLUrl     string      `graphql:"htmlUrl"`
	Owner       EntityOwner `graphql:"owner"`
	Note        string      `graphql:"note"`
}

type DomainConnection struct {
	Nodes      []Domain `json:"nodes"`
	PageInfo   PageInfo `json:"pageInfo"`
	TotalCount int      `json:"totalCount" graphql:"-"`
}

func (d *DomainId) GetTags(client *Client, variables *PayloadVariables) (*TagConnection, error) {
	var q struct {
		Account struct {
			Domain struct {
				Tags TagConnection `graphql:"tags(after: $after, first: $first)"`
			} `graphql:"domain(input: $domain)"`
		}
	}
	if d.Id == "" {
		return nil, fmt.Errorf("Unable to get Tags, invalid domain id: '%s'", d.Id)
	}
	if variables == nil {
		variables = client.InitialPageVariablesPointer()
	}
	(*variables)["domain"] = *NewIdentifier(string(d.Id))

	if err := client.Query(&q, *variables, WithName("DomainTagsList")); err != nil {
		return nil, err
	}
	for q.Account.Domain.Tags.PageInfo.HasNextPage {
		(*variables)["after"] = q.Account.Domain.Tags.PageInfo.End
		resp, err := d.GetTags(client, variables)
		if err != nil {
			return nil, err
		}
		// Add unique tags only
		for _, resp := range resp.Nodes {
			if !slices.Contains[[]Tag, Tag](q.Account.Domain.Tags.Nodes, resp) {
				q.Account.Domain.Tags.Nodes = append(q.Account.Domain.Tags.Nodes, resp)
			}
		}
		q.Account.Domain.Tags.PageInfo = resp.PageInfo
		q.Account.Domain.Tags.TotalCount += resp.TotalCount
	}
	return &q.Account.Domain.Tags, nil
}

func (d *DomainId) ResourceId() ID {
	return d.Id
}

func (d *DomainId) ResourceType() TaggableResource {
	return TaggableResourceDomain
}

func (d *DomainId) ChildSystems(client *Client, variables *PayloadVariables) (*SystemConnection, error) {
	var q struct {
		Account struct {
			Domain struct {
				ChildSystems SystemConnection `graphql:"childSystems(after: $after, first: $first)"`
			} `graphql:"domain(input: $domain)"`
		}
	}
	if d.Id == "" {
		return nil, fmt.Errorf("Unable to get Systems, invalid domain id: '%s'", d.Id)
	}
	if variables == nil {
		variables = client.InitialPageVariablesPointer()
	}

	(*variables)["domain"] = *NewIdentifier(string(d.Id))

	if err := client.Query(&q, *variables, WithName("DomainChildSystemsList")); err != nil {
		return nil, err
	}
	for q.Account.Domain.ChildSystems.PageInfo.HasNextPage {
		(*variables)["after"] = q.Account.Domain.ChildSystems.PageInfo.End
		resp, err := d.ChildSystems(client, variables)
		if err != nil {
			return nil, err
		}
		q.Account.Domain.ChildSystems.Nodes = append(q.Account.Domain.ChildSystems.Nodes, resp.Nodes...)
		q.Account.Domain.ChildSystems.PageInfo = resp.PageInfo
	}
	q.Account.Domain.ChildSystems.TotalCount = len(q.Account.Domain.ChildSystems.Nodes)
	return &q.Account.Domain.ChildSystems, nil
}

func (s *DomainId) AssignSystem(client *Client, systems ...string) error {
	var m struct {
		Payload struct {
			Domain Domain
			Errors []OpsLevelErrors
		} `graphql:"domainChildAssign(domain:$domain, childSystems:$childSystems)"`
	}
	v := PayloadVariables{
		"domain":       *NewIdentifier(string(s.Id)),
		"childSystems": NewIdentifierArray(systems),
	}
	err := client.Mutate(&m, v, WithName("DomainAssignSystem"))
	return HandleErrors(err, m.Payload.Errors)
}

func (c *Client) CreateDomain(input DomainInput) (*Domain, error) {
	var m struct {
		Payload struct {
			Domain Domain
			Errors []OpsLevelErrors
		} `graphql:"domainCreate(input:$input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := c.Mutate(&m, v, WithName("DomainCreate"))
	return &m.Payload.Domain, HandleErrors(err, m.Payload.Errors)
}

func (c *Client) GetDomain(identifier string) (*Domain, error) {
	var q struct {
		Account struct {
			Domain Domain `graphql:"domain(input: $input)"`
		}
	}
	v := PayloadVariables{
		"input": *NewIdentifier(identifier),
	}
	err := c.Query(&q, v, WithName("DomainGet"))
	return &q.Account.Domain, HandleErrors(err, nil)
}

func (c *Client) ListDomains(variables *PayloadVariables) (*DomainConnection, error) {
	var q struct {
		Account struct {
			Domains DomainConnection `graphql:"domains(after: $after, first: $first)"`
		}
	}
	if variables == nil {
		variables = c.InitialPageVariablesPointer()
	}
	if err := c.Query(&q, *variables, WithName("DomainsList")); err != nil {
		return nil, err
	}
	for q.Account.Domains.PageInfo.HasNextPage {
		(*variables)["after"] = q.Account.Domains.PageInfo.End
		resp, err := c.ListDomains(variables)
		if err != nil {
			return nil, err
		}
		q.Account.Domains.Nodes = append(q.Account.Domains.Nodes, resp.Nodes...)
		q.Account.Domains.PageInfo = resp.PageInfo
	}
	q.Account.Domains.TotalCount = len(q.Account.Domains.Nodes)
	return &q.Account.Domains, nil
}

func (c *Client) UpdateDomain(identifier string, input DomainInput) (*Domain, error) {
	var m struct {
		Payload struct {
			Domain Domain
			Errors []OpsLevelErrors
		} `graphql:"domainUpdate(domain:$domain,input:$input)"`
	}
	v := PayloadVariables{
		"domain": *NewIdentifier(identifier),
		"input":  input,
	}
	err := c.Mutate(&m, v, WithName("DomainUpdate"))
	return &m.Payload.Domain, HandleErrors(err, m.Payload.Errors)
}

func (c *Client) DeleteDomain(identifier string) error {
	var d struct {
		Payload struct {
			Errors []OpsLevelErrors `graphql:"errors"`
		} `graphql:"domainDelete(resource: $input)"`
	}
	v := PayloadVariables{
		"input": *NewIdentifier(identifier),
	}
	err := c.Mutate(&d, v, WithName("DomainDelete"))
	return HandleErrors(err, d.Payload.Errors)
}
