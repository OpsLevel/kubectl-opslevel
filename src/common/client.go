package common

import (
	"github.com/opslevel/opslevel-go/v2024"
)

type OpslevelClient struct {
	GetServiceHandler              func(alias string) (*opslevel.Service, error)
	CreateServiceHandler           func(input opslevel.ServiceCreateInput) (*opslevel.Service, error)
	UpdateServiceHandler           func(input opslevel.ServiceUpdateInput) (*opslevel.Service, error)
	CreateAliasHandler             func(input opslevel.AliasCreateInput) error
	AssignTagsHandler              func(service *opslevel.Service, tags map[string]string) error
	AssignPropertyHandler          func(input opslevel.PropertyInput) error
	CreateTagHandler               func(input opslevel.TagCreateInput) error
	CreateToolHandler              func(tool opslevel.ToolCreateInput) error
	GetRepositoryWithAliasHandler  func(alias string) (*opslevel.Repository, error)
	CreateServiceRepositoryHandler func(input opslevel.ServiceRepositoryCreateInput) error
	UpdateServiceRepositoryHandler func(input opslevel.ServiceRepositoryUpdateInput) error
}

func (c *OpslevelClient) GetService(alias string) (*opslevel.Service, error) {
	if c.GetServiceHandler == nil {
		return nil, nil
	}
	return c.GetServiceHandler(alias)
}

func (c *OpslevelClient) CreateService(input opslevel.ServiceCreateInput) (*opslevel.Service, error) {
	if c.CreateServiceHandler == nil {
		return nil, nil
	}
	return c.CreateServiceHandler(input)
}

func (c *OpslevelClient) UpdateService(input opslevel.ServiceUpdateInput) (*opslevel.Service, error) {
	if c.UpdateServiceHandler == nil {
		return nil, nil
	}
	return c.UpdateServiceHandler(input)
}

func (c *OpslevelClient) CreateAlias(input opslevel.AliasCreateInput) error {
	if c.CreateAliasHandler == nil {
		return nil
	}
	return c.CreateAliasHandler(input)
}

func (c *OpslevelClient) AssignTags(service *opslevel.Service, tags map[string]string) error {
	if c.AssignTagsHandler == nil {
		return nil
	}
	return c.AssignTagsHandler(service, tags)
}

func (c *OpslevelClient) AssignProperty(input opslevel.PropertyInput) error {
	if c.AssignPropertyHandler == nil {
		return nil
	}
	return c.AssignPropertyHandler(input)
}

func (c *OpslevelClient) CreateTag(input opslevel.TagCreateInput) error {
	if c.CreateTagHandler == nil {
		return nil
	}
	return c.CreateTagHandler(input)
}

func (c *OpslevelClient) CreateTool(tool opslevel.ToolCreateInput) error {
	if c.CreateToolHandler == nil {
		return nil
	}
	return c.CreateToolHandler(tool)
}

func (c *OpslevelClient) GetRepositoryWithAlias(alias string) (*opslevel.Repository, error) {
	if c.GetRepositoryWithAliasHandler == nil {
		return nil, nil
	}
	return c.GetRepositoryWithAliasHandler(alias)
}

func (c *OpslevelClient) CreateServiceRepository(input opslevel.ServiceRepositoryCreateInput) error {
	if c.CreateServiceRepositoryHandler == nil {
		return nil
	}
	return c.CreateServiceRepositoryHandler(input)
}

func (c *OpslevelClient) UpdateServiceRepository(input opslevel.ServiceRepositoryUpdateInput) error {
	if c.UpdateServiceRepositoryHandler == nil {
		return nil
	}
	return c.UpdateServiceRepositoryHandler(input)
}

func NewOpslevelClient(client *opslevel.Client) *OpslevelClient {
	return &OpslevelClient{
		GetServiceHandler: func(alias string) (*opslevel.Service, error) {
			return client.GetServiceWithAlias(alias)
		},
		CreateServiceHandler: func(input opslevel.ServiceCreateInput) (*opslevel.Service, error) {
			return client.CreateService(input)
		},
		UpdateServiceHandler: func(input opslevel.ServiceUpdateInput) (*opslevel.Service, error) {
			return client.UpdateService(opslevel.ConvertServiceUpdateInput(input))
		},
		CreateAliasHandler: func(input opslevel.AliasCreateInput) error {
			_, err := client.CreateAlias(input)
			return err
		},
		AssignTagsHandler: func(service *opslevel.Service, tags map[string]string) error {
			_, err := client.AssignTags(string(service.Id), tags)
			return err
		},
		AssignPropertyHandler: func(input opslevel.PropertyInput) error {
			_, err := client.PropertyAssign(input)
			return err
		},
		CreateTagHandler: func(input opslevel.TagCreateInput) error {
			_, err := client.CreateTag(input)
			return err
		},
		CreateToolHandler: func(tool opslevel.ToolCreateInput) error {
			_, err := client.CreateTool(tool)
			return err
		},
		GetRepositoryWithAliasHandler: func(alias string) (*opslevel.Repository, error) {
			return client.GetRepositoryWithAlias(alias)
		},
		CreateServiceRepositoryHandler: func(input opslevel.ServiceRepositoryCreateInput) error {
			_, err := client.CreateServiceRepository(input)
			return err
		},
		UpdateServiceRepositoryHandler: func(input opslevel.ServiceRepositoryUpdateInput) error {
			_, err := client.UpdateServiceRepository(input)
			return err
		},
	}
}
