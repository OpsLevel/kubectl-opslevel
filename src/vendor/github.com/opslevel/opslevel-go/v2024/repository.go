package opslevel

import (
	"fmt"
	"slices"

	"github.com/relvacode/iso8601"
)

type Language struct {
	Name  string
	Usage float32
}

// Lightweight Repository struct used to make some API calls return less data
type RepositoryId struct {
	Id           ID
	DefaultAlias string
}

type Repository struct {
	ArchivedAt         iso8601.Time
	CreatedOn          iso8601.Time
	DefaultAlias       string
	DefaultBranch      string
	Description        string
	Forked             bool
	HtmlUrl            string
	Id                 ID
	Languages          []Language
	LastOwnerChangedAt iso8601.Time
	Name               string
	Organization       string
	Owner              TeamId
	Private            bool
	RepoKey            string
	Services           *RepositoryServiceConnection
	Tags               *TagConnection
	Tier               Tier
	Type               string
	Url                string
	Visible            bool
}

type RepositoryPath struct {
	Href string
	Path string
}

type ServiceRepository struct {
	BaseDirectory string
	DisplayName   string
	Id            ID
	Repository    RepositoryId
	Service       ServiceId
}

type RepositoryConnection struct {
	HiddenCount       int
	Nodes             []Repository
	OrganizationCount int
	OwnedCount        int
	PageInfo          PageInfo
	TotalCount        int
	VisibleCount      int
}

type RepositoryServiceEdge struct {
	AtRoot              bool
	Node                ServiceId
	Paths               []RepositoryPath
	ServiceRepositories []ServiceRepository
}

type RepositoryServiceConnection struct {
	Edges      []RepositoryServiceEdge
	PageInfo   PageInfo
	TotalCount int
}

type ServiceRepositoryEdge struct {
	Node                RepositoryId
	ServiceRepositories []ServiceRepository
}

type ServiceRepositoryConnection struct {
	Edges      []ServiceRepositoryEdge
	PageInfo   PageInfo
	TotalCount int
}

func (r *Repository) ResourceId() ID {
	return r.Id
}

func (r *Repository) ResourceType() TaggableResource {
	return TaggableResourceRepository
}

func (r *Repository) GetService(service ID, directory string) *ServiceRepository {
	for _, edge := range r.Services.Edges {
		for _, connection := range edge.ServiceRepositories {
			if connection.Service.Id == service && connection.BaseDirectory == directory {
				return &connection
			}
		}
	}
	return nil
}

func (r *Repository) Hydrate(client *Client) error {
	if r.Services == nil {
		r.Services = &RepositoryServiceConnection{}
	}
	if r.Services.PageInfo.HasNextPage {
		variables := client.InitialPageVariablesPointer()
		(*variables)["after"] = r.Services.PageInfo.End
		_, err := r.GetServices(client, variables)
		if err != nil {
			return err
		}
	}

	if r.Tags == nil {
		r.Tags = &TagConnection{}
	}
	if r.Tags.PageInfo.HasNextPage {
		variables := client.InitialPageVariablesPointer()
		(*variables)["after"] = r.Tags.PageInfo.End
		_, err := r.GetTags(client, variables)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) GetServices(client *Client, variables *PayloadVariables) (*RepositoryServiceConnection, error) {
	var q struct {
		Account struct {
			Repository struct {
				Services RepositoryServiceConnection `graphql:"services(after: $after, first: $first)"`
			} `graphql:"repository(id: $id)"`
		}
	}
	if r.Id == "" {
		return nil, fmt.Errorf("Unable to get Services, invalid repository id: '%s'", r.Id)
	}
	if variables == nil {
		variables = client.InitialPageVariablesPointer()
	}
	(*variables)["id"] = r.Id
	if err := client.Query(&q, *variables, WithName("RepositoryServicesList")); err != nil {
		return nil, err
	}
	if r.Services == nil {
		r.Services = &RepositoryServiceConnection{}
	}
	r.Services.Edges = append(r.Services.Edges, q.Account.Repository.Services.Edges...)
	r.Services.PageInfo = q.Account.Repository.Services.PageInfo
	r.Services.TotalCount += q.Account.Repository.Services.TotalCount
	for r.Services.PageInfo.HasNextPage {
		(*variables)["after"] = r.Services.PageInfo.End
		_, err := r.GetServices(client, variables)
		if err != nil {
			return nil, err
		}
	}
	return r.Services, nil
}

func (r *Repository) GetTags(client *Client, variables *PayloadVariables) (*TagConnection, error) {
	var q struct {
		Account struct {
			Repository struct {
				Tags TagConnection `graphql:"tags(after: $after, first: $first)"`
			} `graphql:"repository(id: $id)"`
		}
	}
	if r.Id == "" {
		return nil, fmt.Errorf("Unable to get Tags, invalid repository id: '%s'", r.Id)
	}
	if variables == nil {
		variables = client.InitialPageVariablesPointer()
	}
	(*variables)["id"] = r.Id
	if err := client.Query(&q, *variables, WithName("RepositoryTagsList")); err != nil {
		return nil, err
	}
	if r.Tags == nil {
		r.Tags = &TagConnection{}
	}
	// Add unique tags only
	for _, tagNode := range q.Account.Repository.Tags.Nodes {
		if !slices.Contains[[]Tag, Tag](r.Tags.Nodes, tagNode) {
			r.Tags.Nodes = append(r.Tags.Nodes, tagNode)
		}
	}
	r.Tags.PageInfo = q.Account.Repository.Tags.PageInfo
	r.Tags.TotalCount += q.Account.Repository.Tags.TotalCount
	for r.Tags.PageInfo.HasNextPage {
		(*variables)["after"] = r.Tags.PageInfo.End
		_, err := r.GetTags(client, variables)
		if err != nil {
			return nil, err
		}
	}
	return r.Tags, nil
}

//#region Create

func (client *Client) ConnectServiceRepository(service *ServiceId, repository *Repository) (*ServiceRepository, error) {
	input := ServiceRepositoryCreateInput{
		Service:       *NewIdentifier(string(service.Id)),
		Repository:    *NewIdentifier(string(repository.Id)),
		BaseDirectory: RefOf("/"),
		DisplayName:   RefOf(fmt.Sprintf("%s/%s", repository.Organization, repository.Name)),
	}
	return client.CreateServiceRepository(input)
}

func (client *Client) CreateServiceRepository(input ServiceRepositoryCreateInput) (*ServiceRepository, error) {
	var m struct {
		Payload struct {
			ServiceRepository ServiceRepository
			Errors            []OpsLevelErrors
		} `graphql:"serviceRepositoryCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("ServiceRepositoryCreate"))
	return &m.Payload.ServiceRepository, HandleErrors(err, m.Payload.Errors)
}

//#endregion

//#region Retrieve

func (client *Client) GetRepositoryWithAlias(alias string) (*Repository, error) {
	var q struct {
		Account struct {
			Repository Repository `graphql:"repository(alias: $repo)"`
		}
	}
	v := PayloadVariables{
		"repo": alias,
	}
	if err := client.Query(&q, v, WithName("RepositoryGet")); err != nil {
		return nil, err
	}
	if err := q.Account.Repository.Hydrate(client); err != nil {
		return &q.Account.Repository, err
	}
	return &q.Account.Repository, nil
}

func (client *Client) GetRepository(id ID) (*Repository, error) {
	var q struct {
		Account struct {
			Repository Repository `graphql:"repository(id: $repo)"`
		}
	}
	v := PayloadVariables{
		"repo": id,
	}
	if err := client.Query(&q, v, WithName("RepositoryGet")); err != nil {
		return nil, err
	}
	if err := q.Account.Repository.Hydrate(client); err != nil {
		return &q.Account.Repository, err
	}
	return &q.Account.Repository, nil
}

func (client *Client) ListRepositories(variables *PayloadVariables) (*RepositoryConnection, error) {
	var q struct {
		Account struct {
			Repositories RepositoryConnection `graphql:"repositories(after: $after, first: $first)"`
		}
	}
	if variables == nil {
		variables = client.InitialPageVariablesPointer()
	}
	if err := client.Query(&q, *variables, WithName("RepositoryList")); err != nil {
		return &q.Account.Repositories, err
	}
	for q.Account.Repositories.PageInfo.HasNextPage {
		(*variables)["after"] = q.Account.Repositories.PageInfo.End
		resp, err := client.ListRepositories(variables)
		if err != nil {
			return nil, err
		}
		for _, node := range resp.Nodes {
			err := node.Hydrate(client)
			if err != nil {
				return nil, err
			}
			q.Account.Repositories.Nodes = append(q.Account.Repositories.Nodes, node)
		}
		q.Account.Repositories.PageInfo = resp.PageInfo
		q.Account.Repositories.TotalCount += resp.TotalCount
	}
	return &q.Account.Repositories, nil
}

func (client *Client) ListRepositoriesWithTier(tier string, variables *PayloadVariables) (*RepositoryConnection, error) {
	var q struct {
		Account struct {
			Repositories RepositoryConnection `graphql:"repositories(tierAlias: $tier, after: $after, first: $first)"`
		}
	}
	if variables == nil {
		variables = client.InitialPageVariablesPointer()
	}
	(*variables)["tier"] = tier
	if err := client.Query(&q, *variables, WithName("RepositoryListWithTier")); err != nil {
		return &q.Account.Repositories, err
	}
	for q.Account.Repositories.PageInfo.HasNextPage {
		(*variables)["after"] = q.Account.Repositories.PageInfo.End
		resp, err := client.ListRepositoriesWithTier(tier, variables)
		if err != nil {
			return nil, err
		}
		for _, node := range resp.Nodes {
			err := node.Hydrate(client)
			if err != nil {
				return nil, err
			}
			q.Account.Repositories.Nodes = append(q.Account.Repositories.Nodes, node)
		}
		q.Account.Repositories.PageInfo = resp.PageInfo
		q.Account.Repositories.TotalCount += resp.TotalCount
	}
	return &q.Account.Repositories, nil
}

//#endregion

//#region Update

func (client *Client) UpdateRepository(input RepositoryUpdateInput) (*Repository, error) {
	var m struct {
		Payload struct {
			Repository Repository
			Errors     []OpsLevelErrors
		} `graphql:"repositoryUpdate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("RepositoryUpdate"))
	return &m.Payload.Repository, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) UpdateServiceRepository(input ServiceRepositoryUpdateInput) (*ServiceRepository, error) {
	var m struct {
		Payload struct {
			ServiceRepository ServiceRepository
			Errors            []OpsLevelErrors
		} `graphql:"serviceRepositoryUpdate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("ServiceRepositoryUpdate"))
	return &m.Payload.ServiceRepository, HandleErrors(err, m.Payload.Errors)
}

//#endregion

//#region Delete

func (client *Client) DeleteServiceRepository(id ID) error {
	var m struct {
		Payload struct {
			Id     ID `graphql:"deletedId"`
			Errors []OpsLevelErrors
		} `graphql:"serviceRepositoryDelete(input: $input)"`
	}
	v := PayloadVariables{
		"input": DeleteInput{Id: id},
	}
	err := client.Mutate(&m, v, WithName("ServiceRepositoryDelete"))
	return HandleErrors(err, m.Payload.Errors)
}

//#endregion
