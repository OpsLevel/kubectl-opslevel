package common_test

import (
	"fmt"

	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/opslevel/opslevel-go/v2024"
)

type TestClientBuilder struct {
	client common.OpslevelClient
}

// NewTestClientBuilderWithError creates a client where every function that has not been explicitly set will error.
func NewTestClientBuilderWithError() *TestClientBuilder {
	return &TestClientBuilder{
		client: common.OpslevelClient{
			GetServiceHandler: func(alias string) (*opslevel.Service, error) {
				return nil, fmt.Errorf("should not be called")
			},
			CreateServiceHandler: func(input opslevel.ServiceCreateInput) (*opslevel.Service, error) {
				return nil, fmt.Errorf("should not be called")
			},
			UpdateServiceHandler: func(input opslevel.ServiceUpdateInput) (*opslevel.Service, error) {
				return nil, fmt.Errorf("should not be called")
			},
			CreateAliasHandler: func(input opslevel.AliasCreateInput) error {
				return fmt.Errorf("should not be called")
			},
			AssignTagsHandler: func(service *opslevel.Service, tags map[string]string) error {
				return fmt.Errorf("should not be called")
			},
			AssignPropertyHandler: func(input opslevel.PropertyInput) error {
				return fmt.Errorf("should not be called")
			},
			CreateTagHandler: func(input opslevel.TagCreateInput) error {
				return fmt.Errorf("should not be called")
			},
			CreateToolHandler: func(tool opslevel.ToolCreateInput) error {
				return fmt.Errorf("should not be called")
			},
			GetRepositoryWithAliasHandler: func(alias string) (*opslevel.Repository, error) {
				return nil, fmt.Errorf("should not be called")
			},
			CreateServiceRepositoryHandler: func(input opslevel.ServiceRepositoryCreateInput) error {
				return fmt.Errorf("should not be called")
			},
			UpdateServiceRepositoryHandler: func(input opslevel.ServiceRepositoryUpdateInput) error {
				return fmt.Errorf("should not be called")
			},
		},
	}
}

// NewTestClientBuilderWithPanic creates a client where every function that has not been explicitly set will panic.
func NewTestClientBuilderWithPanic() *TestClientBuilder {
	return &TestClientBuilder{
		client: common.OpslevelClient{
			GetServiceHandler: func(alias string) (*opslevel.Service, error) {
				panic("should not be called")
			},
			CreateServiceHandler: func(input opslevel.ServiceCreateInput) (*opslevel.Service, error) {
				panic("should not be called")
			},
			UpdateServiceHandler: func(input opslevel.ServiceUpdateInput) (*opslevel.Service, error) {
				panic("should not be called")
			},
			CreateAliasHandler: func(input opslevel.AliasCreateInput) error {
				panic("should not be called")
			},
			AssignTagsHandler: func(service *opslevel.Service, tags map[string]string) error {
				panic("should not be called")
			},
			AssignPropertyHandler: func(input opslevel.PropertyInput) error {
				panic("should not be called")
			},
			CreateTagHandler: func(input opslevel.TagCreateInput) error {
				panic("should not be called")
			},
			CreateToolHandler: func(tool opslevel.ToolCreateInput) error {
				panic("should not be called")
			},
			GetRepositoryWithAliasHandler: func(alias string) (*opslevel.Repository, error) {
				panic("should not be called")
			},
			CreateServiceRepositoryHandler: func(input opslevel.ServiceRepositoryCreateInput) error {
				panic("should not be called")
			},
			UpdateServiceRepositoryHandler: func(input opslevel.ServiceRepositoryUpdateInput) error {
				panic("should not be called")
			},
		},
	}
}

func (b *TestClientBuilder) SetAssignPropertyHandler(fn func(input opslevel.PropertyInput) error) *TestClientBuilder {
	b.client.AssignPropertyHandler = fn
	return b
}

func (b *TestClientBuilder) SetAssignTagsHandler(fn func(service *opslevel.Service, tags map[string]string) error) *TestClientBuilder {
	b.client.AssignTagsHandler = fn
	return b
}

func (b *TestClientBuilder) SetCreateAliasHandler(fn func(input opslevel.AliasCreateInput) error) *TestClientBuilder {
	b.client.CreateAliasHandler = fn
	return b
}

func (b *TestClientBuilder) SetCreateServiceHandler(fn func(input opslevel.ServiceCreateInput) (*opslevel.Service, error)) *TestClientBuilder {
	b.client.CreateServiceHandler = fn
	return b
}

func (b *TestClientBuilder) SetCreateServiceRepositoryHandler(fn func(input opslevel.ServiceRepositoryCreateInput) error) *TestClientBuilder {
	b.client.CreateServiceRepositoryHandler = fn
	return b
}

func (b *TestClientBuilder) SetCreateTagHandler(fn func(input opslevel.TagCreateInput) error) *TestClientBuilder {
	b.client.CreateTagHandler = fn
	return b
}

func (b *TestClientBuilder) SetCreateToolHandler(fn func(input opslevel.ToolCreateInput) error) *TestClientBuilder {
	b.client.CreateToolHandler = fn
	return b
}

func (b *TestClientBuilder) SetGetRepositoryWithAliasHandler(fn func(alias string) (*opslevel.Repository, error)) *TestClientBuilder {
	b.client.GetRepositoryWithAliasHandler = fn
	return b
}

func (b *TestClientBuilder) SetGetServiceHandler(fn func(alias string) (*opslevel.Service, error)) *TestClientBuilder {
	b.client.GetServiceHandler = fn
	return b
}

func (b *TestClientBuilder) SetUpdateServiceHandler(fn func(input opslevel.ServiceUpdateInput) (*opslevel.Service, error)) *TestClientBuilder {
	b.client.UpdateServiceHandler = fn
	return b
}

func (b *TestClientBuilder) SetUpdateServiceRepositoryHandler(fn func(input opslevel.ServiceRepositoryUpdateInput) error) *TestClientBuilder {
	b.client.UpdateServiceRepositoryHandler = fn
	return b
}

func (b *TestClientBuilder) GetClient() *common.OpslevelClient {
	return &b.client
}
