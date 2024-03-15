package common_test

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/opslevel/opslevel-go/v2024"
	opslevel_jq_parser "github.com/opslevel/opslevel-jq-parser/v2024"
	"github.com/rocktavious/autopilot/v2023"
)

var gotService = false

type TestClientBuilder struct {
	client *common.OpslevelClient
}

// TODO: add testing for selector

//go:embed testdata/full_registration.json
var fullRegistrationJSON []byte
var fullRegistration opslevel_jq_parser.ServiceRegistration
var fullRegistrationCreateInput opslevel.ServiceCreateInput
var fullRegistrationUpdateInput opslevel.ServiceUpdateInput

//go:embed testdata/min_registration.json
var minRegistrationJSON []byte
var minRegistration opslevel_jq_parser.ServiceRegistration
var minRegistrationCreateInput opslevel.ServiceCreateInput
var minRegistrationUpdateInput opslevel.ServiceUpdateInput

func service(registration opslevel_jq_parser.ServiceRegistration) *opslevel.Service {
	svc := opslevel.Service{
		Name: registration.Name,
	}
	if registration.Description != "" {
		svc.Description = registration.Description
	}
	if registration.Framework != "" {
		svc.Framework = registration.Framework
	}
	if registration.Language != "" {
		svc.Language = registration.Language
	}
	if registration.Lifecycle != "" {
		svc.Lifecycle = opslevel.Lifecycle{Alias: registration.Lifecycle}
	}
	if registration.Owner != "" {
		svc.Owner = opslevel.TeamId{Alias: registration.Owner}
	}
	if registration.System != "" {
		svc.Parent = &opslevel.SystemId{Aliases: []string{registration.System}}
	}
	if registration.Product != "" {
		svc.Product = registration.Product
	}
	if registration.Tier != "" {
		svc.Tier = opslevel.Tier{Alias: registration.Tier}
	}
	svc.Id = "XXX"
	return &svc
}

func createInput(registration opslevel_jq_parser.ServiceRegistration) opslevel.ServiceCreateInput {
	input := opslevel.ServiceCreateInput{
		Name: registration.Name,
	}
	if registration.Description != "" {
		input.Description = opslevel.RefOf(registration.Description)
	}
	if registration.Framework != "" {
		input.Framework = opslevel.RefOf(registration.Framework)
	}
	if registration.Language != "" {
		input.Language = opslevel.RefOf(registration.Language)
	}
	if registration.Lifecycle != "" {
		input.LifecycleAlias = opslevel.RefOf(registration.Lifecycle)
	}
	if registration.Owner != "" {
		input.OwnerInput = opslevel.NewIdentifier(registration.Owner)
	}
	if registration.System != "" {
		input.Parent = opslevel.NewIdentifier(registration.System)
	}
	if registration.Product != "" {
		input.Product = opslevel.RefOf(registration.Product)
	}
	if registration.Tier != "" {
		input.TierAlias = opslevel.RefOf(registration.Tier)
	}
	return input
}

func updateInput(registration opslevel_jq_parser.ServiceRegistration) opslevel.ServiceUpdateInput {
	serviceInput := opslevel.ServiceUpdateInput{
		Id: opslevel.NewID("XXX"),
	}
	if registration.Description != "" {
		serviceInput.Description = opslevel.RefOf(registration.Description)
	}
	if registration.Framework != "" {
		serviceInput.Framework = opslevel.RefOf(registration.Framework)
	}
	if registration.Language != "" {
		serviceInput.Language = opslevel.RefOf(registration.Language)
	}
	if registration.Lifecycle != "" {
		serviceInput.LifecycleAlias = opslevel.RefOf(registration.Lifecycle)
	}
	if registration.Name != "" {
		serviceInput.Name = opslevel.RefOf(registration.Name)
	}
	if registration.Owner != "" {
		serviceInput.OwnerInput = opslevel.NewIdentifier(registration.Owner)
	}
	if registration.System != "" {
		serviceInput.Parent = opslevel.NewIdentifier(registration.System)
	}
	if registration.Product != "" {
		serviceInput.Product = opslevel.RefOf(registration.Product)
	}
	if registration.Tier != "" {
		serviceInput.TierAlias = opslevel.RefOf(registration.Tier)
	}
	return serviceInput
}

// initialize embedded test data
func init() {
	var err error
	err = json.Unmarshal(fullRegistrationJSON, &fullRegistration)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(minRegistrationJSON, &minRegistration)
	if err != nil {
		panic(err)
	}

	// expected inputs
	fullRegistrationCreateInput = createInput(fullRegistration)
	fullRegistrationUpdateInput = updateInput(fullRegistration)
	minRegistrationCreateInput = createInput(minRegistration)
	minRegistrationUpdateInput = updateInput(minRegistration)
}

// NewTestClientBuilder creates a client where every function that has not been explicitly set will panic. Useful for testing
// what path the reconciler will take, and also for exiting early by returning an error that can be added to an expectation.
func NewTestClientBuilder() *TestClientBuilder {
	return &TestClientBuilder{
		client: &common.OpslevelClient{
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
	return b.client
}

func TestReconcilerReconcile(t *testing.T) {
	// Arrange
	type TestCase struct {
		assert                 func(t *testing.T, err error)
		client                 *common.OpslevelClient
		disableServiceCreation bool
		name                   string
		registration           opslevel_jq_parser.ServiceRegistration
	}

	cases := []TestCase{
		{
			name:         "Missing Name Should Error",
			registration: opslevel_jq_parser.ServiceRegistration{},
			client:       NewTestClientBuilder().GetClient(),
			assert: func(t *testing.T, err error) {
				autopilot.Equals(t, "cannot reconcile service with no name", err.Error())
			},
		},
		{
			name: "Missing Aliases Should Error",
			registration: opslevel_jq_parser.ServiceRegistration{
				Name: "test",
			},
			client: NewTestClientBuilder().GetClient(),
			assert: func(t *testing.T, err error) {
				autopilot.Equals(t, "[test] found 0 aliases from kubernetes data", err.Error())
			},
		},
		{
			name: "Multiple Matching Aliases Should Halt",
			registration: opslevel_jq_parser.ServiceRegistration{
				Aliases: []string{"test1", "test2", "test3"},
				Name:    "test",
			},
			client: NewTestClientBuilder().SetGetServiceHandler(func(alias string) (*opslevel.Service, error) {
				return &opslevel.Service{}, nil
			}).GetClient(),
			assert: func(t *testing.T, err error) {
				autopilot.Equals(t, "[test] found multiple services with aliases = [[test1 test2 test3]]. cannot know which service to target for update ... skipping reconciliation", err.Error())
			},
		},
		{
			name: "API Error On Get Service Should Halt",
			registration: opslevel_jq_parser.ServiceRegistration{
				Aliases: []string{"test"},
				Name:    "test",
			},
			client: NewTestClientBuilder().SetGetServiceHandler(func(alias string) (*opslevel.Service, error) {
				return nil, fmt.Errorf("api error")
			}).GetClient(),
			assert: func(t *testing.T, err error) {
				autopilot.Equals(t, "[test] api error during service lookup by alias.  unable to guarantee service was found or not ... skipping reconciliation", err.Error())
			},
		},
		{
			name: "API Error On Create Service Should Halt",
			registration: opslevel_jq_parser.ServiceRegistration{
				Aliases: []string{"test"},
				Name:    "test",
			},
			client: NewTestClientBuilder().SetGetServiceHandler(func(alias string) (*opslevel.Service, error) {
				return nil, nil
			}).SetCreateServiceHandler(func(input opslevel.ServiceCreateInput) (*opslevel.Service, error) {
				return nil, fmt.Errorf("api error from create")
			}).GetClient(),
			assert: func(t *testing.T, err error) {
				autopilot.Equals(t, "api error from create", err.Error())
			},
		},
		{
			name: "API Error On Update Service Should Halt",
			registration: opslevel_jq_parser.ServiceRegistration{
				Aliases:     []string{"test"},
				Name:        "test",
				Description: "testing 123",
			},
			client: NewTestClientBuilder().SetGetServiceHandler(func(alias string) (*opslevel.Service, error) {
				return &opslevel.Service{Name: "test"}, nil
			}).SetCreateServiceHandler(func(input opslevel.ServiceCreateInput) (*opslevel.Service, error) {
				return &opslevel.Service{Name: "test"}, nil
			}).SetUpdateServiceHandler(func(input opslevel.ServiceUpdateInput) (*opslevel.Service, error) {
				return nil, fmt.Errorf("api error from update")
			}).GetClient(),
			assert: func(t *testing.T, err error) {
				autopilot.Equals(t, "api error from update", err.Error())
			},
		},
		{
			name: "Happy Path Do Not Create Services",
			registration: opslevel_jq_parser.ServiceRegistration{
				Aliases: []string{"test"},
				Name:    "test",
			},
			client: NewTestClientBuilder().SetGetServiceHandler(func(alias string) (*opslevel.Service, error) {
				return nil, nil
			}).GetClient(),
			assert: func(t *testing.T, err error) {
				autopilot.Ok(t, err)
			},
			disableServiceCreation: true,
		},
		// TODO: tags assign case: service tags not populated - tags are populated
		// TODO: tags create case: service tags not populated - tags are populated
		// TODO: tools - have none/have 1/have 2 but 1 needs to be updated
		// TODO: repositories - have none/have 1/have 2 but 1 needs to be updated
		// TODO: properties - have none/have 1/have 2 but 1 needs to be updated
		{
			name:         "Happy Path Create Service Continue Executing On Error After Service Update",
			registration: fullRegistration,
			client: NewTestClientBuilder().SetGetServiceHandler(func(alias string) (*opslevel.Service, error) {
				if gotService == true {
					return nil, nil
				}
				gotService = true
				return service(fullRegistration), nil
			}).SetUpdateServiceHandler(func(input opslevel.ServiceUpdateInput) (*opslevel.Service, error) {
				return service(fullRegistration), nil
			}).SetCreateAliasHandler(func(input opslevel.AliasCreateInput) error {
				return fmt.Errorf("create alias error")
			}).SetAssignTagsHandler(func(service *opslevel.Service, tags map[string]string) error {
				return fmt.Errorf("assign tags error")
			}).SetCreateTagHandler(func(input opslevel.TagCreateInput) error {
				return fmt.Errorf("create tag error")
			}).SetCreateToolHandler(func(input opslevel.ToolCreateInput) error {
				return fmt.Errorf("create tool error")
			}).SetGetRepositoryWithAliasHandler(func(alias string) (*opslevel.Repository, error) {
				return nil, fmt.Errorf("get repo error")
			}).SetCreateServiceRepositoryHandler(func(input opslevel.ServiceRepositoryCreateInput) error {
				return fmt.Errorf("create repo error")
			}).SetAssignPropertyHandler(func(input opslevel.PropertyInput) error {
				return fmt.Errorf("create property error")
			}).GetClient(),
			assert: func(t *testing.T, err error) {
				// TODO: how can we ensure that every single func was called
				autopilot.Ok(t, err)
			},
		},
		{
			name: "Happy Path Do Everything",
			// TODO: happy path with update
			// TODO: with min registration (use update)
			// TODO: validate steps for assign/create tag, create tool, create repo, create property
			registration: fullRegistration,
			client: NewTestClientBuilder().SetGetServiceHandler(func(alias string) (*opslevel.Service, error) {
				return nil, nil
			}).SetCreateServiceHandler(func(input opslevel.ServiceCreateInput) (*opslevel.Service, error) {
				return service(fullRegistration), nil
			}).SetCreateAliasHandler(func(input opslevel.AliasCreateInput) error {
				return nil
			}).SetAssignTagsHandler(func(service *opslevel.Service, tags map[string]string) error {
				// TODO: the service.Tags is nil.
				return nil
			}).SetCreateTagHandler(func(input opslevel.TagCreateInput) error {
				return nil
			}).SetCreateToolHandler(func(input opslevel.ToolCreateInput) error {
				return nil
			}).SetGetRepositoryWithAliasHandler(func(alias string) (*opslevel.Repository, error) {
				return nil, nil
			}).SetCreateServiceRepositoryHandler(func(input opslevel.ServiceRepositoryCreateInput) error {
				return nil
			}).SetAssignPropertyHandler(func(input opslevel.PropertyInput) error {
				return nil
			}).GetClient(),
			assert: func(t *testing.T, err error) {
				autopilot.Ok(t, err)
			},
		},
	}
	// Act
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// TODO: this is a hack to let service updates work by making the get service call only return once if there are multiple aliases that can return multiple services.
			switch tc.name {
			case "Happy Path Create Service Continue Executing On Error After Service Update":
				gotService = false
			}
			result := common.NewServiceReconciler(tc.client, tc.disableServiceCreation).Reconcile(tc.registration)
			tc.assert(t, result)
		})
	}
}

func Test_Reconciler_ContainsAllTags(t *testing.T) {
	// Arrange
	type TestCase struct {
		input  []opslevel.TagInput
		tags   []opslevel.Tag
		result bool
	}
	reconciler := common.NewServiceReconciler(&common.OpslevelClient{}, false)
	cases := map[string]TestCase{
		"Is True When All Tags Overlap": {
			input: []opslevel.TagInput{
				{
					Key:   "foo",
					Value: "bar",
				},
				{
					Key:   "hello",
					Value: "world",
				},
				{
					Key:   "apple",
					Value: "orange",
				},
			},
			tags: []opslevel.Tag{
				{
					Key:   "foo",
					Value: "bar",
				},
				{
					Key:   "hello",
					Value: "world",
				},
				{
					Key:   "apple",
					Value: "orange",
				},
			},
			result: true,
		},
		"Is False When Tag Input Has More": {
			input: []opslevel.TagInput{
				{
					Key:   "foo",
					Value: "bar",
				},
				{
					Key:   "hello",
					Value: "world",
				},
				{
					Key:   "apple",
					Value: "orange",
				},
			},
			tags: []opslevel.Tag{
				{
					Key:   "foo",
					Value: "bar",
				},
				{
					Key:   "hello",
					Value: "world",
				},
			},
			result: false,
		},
		"Is True When Service Has More": {
			input: []opslevel.TagInput{
				{
					Key:   "foo",
					Value: "bar",
				},
			},
			tags: []opslevel.Tag{
				{
					Key:   "foo",
					Value: "bar",
				},
				{
					Key:   "hello",
					Value: "world",
				},
				{
					Key:   "apple",
					Value: "orange",
				},
			},
			result: true,
		},
	}
	// Act
	autopilot.RunTableTests(t, cases, func(t *testing.T, test TestCase) {
		// Assert
		autopilot.Equals(t, test.result, reconciler.ContainsAllTags(test.input, test.tags))
	})
}

func Test_Reconciler_HandleTools(t *testing.T) {
	// Arrange
	registration := opslevel_jq_parser.ServiceRegistration{
		Aliases: []string{"a_test_service"},
		Name:    "A test service",
		Tools:   newToolInputs("A", "B", "C", "D", "E", "F", "G"),
	}
	service := opslevel.Service{
		ServiceId: opslevel.ServiceId{
			Id: opslevel.ID("XXX"),
		},
		Name: "ATestService",
		Tools: &opslevel.ToolConnection{
			Nodes: newTools("A", "B", "C", "D", "E"),
		},
	}
	toolsCreated := make([]opslevel.ToolCreateInput, 0)
	reconciler := common.NewServiceReconciler(&common.OpslevelClient{
		GetServiceHandler: func(alias string) (*opslevel.Service, error) {
			return &service, nil
		},
		CreateToolHandler: func(tool opslevel.ToolCreateInput) error {
			toolsCreated = append(toolsCreated, tool)
			return nil
		},
	}, false)
	// Act
	err := reconciler.Reconcile(registration)
	autopilot.Ok(t, err)
	// Assert
	autopilot.Assert(t, len(toolsCreated) == 2 && toolsCreated[0].DisplayName == "F" &&
		toolsCreated[1].DisplayName == "G", "expected tools created to be ['F','G']")
}

func Test_Reconciler_HandleProperties(t *testing.T) {
	// Arrange
	props := []opslevel.PropertyInput{
		{
			Definition: *opslevel.NewIdentifier("prop_bool"),
			Value:      opslevel.JsonString("true"),
		},
		{
			Definition: *opslevel.NewIdentifier("prop_empty_object"),
			Value:      opslevel.JsonString("{"),
		},
		{
			Definition: *opslevel.NewIdentifier("prop_object"),
			Value:      opslevel.JsonString("{\"message\":\"hello world\",\"condition\":true}"),
		},
		{
			Definition: *opslevel.NewIdentifier("prop_string"),
			Value:      opslevel.JsonString("hello world"),
		},
	}
	registration := opslevel_jq_parser.ServiceRegistration{
		Aliases:    []string{"a_test_service_with_properties"},
		Name:       "A test service with properties",
		Properties: props,
	}
	service := opslevel.Service{
		ServiceId: opslevel.ServiceId{
			Id: opslevel.ID("XXX"),
		},
		Name:       "ATestServiceWithProperties",
		Properties: nil,
	}
	results := make([]opslevel.PropertyInput, 0)
	reconciler := common.NewServiceReconciler(&common.OpslevelClient{
		GetServiceHandler: func(alias string) (*opslevel.Service, error) {
			return &service, nil
		},
		AssignPropertyHandler: func(input opslevel.PropertyInput) error {
			results = append(results, input)
			return nil
		},
	}, false)

	// Act
	err := reconciler.Reconcile(registration)
	autopilot.Ok(t, err)

	// Assert
	expLen, gotLen := len(props), len(results)
	autopilot.Assert(t, gotLen == expLen, fmt.Sprintf("expected to get %d property assignments got %d", expLen, gotLen))
	// TODO: validate the results are correct
	//for _, x := range results {
	//	def := *x.Definition.Alias
	//	expId, gotId := service.ServiceId.Id, *x.Owner.Id
	//	autopilot.Assert(t, gotId == expId, fmt.Sprintf("[%s] unexpected owner ID '%s' - does not match service ID '%s'", def, gotId, expId))
	//
	//	if _, ok := props[def]; ok {
	//		value, err := opslevel.NewJSONInput(val)
	//		if err != nil {
	//			log.Error().Err(err).Msgf("[%s] Failed parsing property: '%s'", service.Name, def)
	//			continue
	//		}
	//		expVal := string(*value)
	//		gotVal := string(x.Value)
	//		autopilot.Assert(t, gotVal == expVal, fmt.Sprintf("[%s] expected value for to be: '%s' got: '%s'", def, expVal, gotVal))
	//	} else {
	//		autopilot.Ok(t, fmt.Errorf("unexpected property definition alias: '%s'", def))
	//	}
	//}
}

func newToolInputs(names ...string) []opslevel.ToolCreateInput {
	inputs := make([]opslevel.ToolCreateInput, len(names))
	for i, d := range names {
		inputs[i] = opslevel.ToolCreateInput{
			Category:    opslevel.ToolCategoryCode,
			DisplayName: d,
			Environment: opslevel.RefOf("XXX"),
		}
	}
	return inputs
}

func newTools(names ...string) []opslevel.Tool {
	tools := make([]opslevel.Tool, len(names))
	for i, d := range names {
		tools[i] = opslevel.Tool{
			Category:    opslevel.ToolCategoryCode,
			DisplayName: d,
			Environment: "XXX",
		}
	}
	return tools
}
