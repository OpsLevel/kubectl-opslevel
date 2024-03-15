package common_test

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/rs/zerolog/log"

	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/opslevel/opslevel-go/v2024"
	opslevel_jq_parser "github.com/opslevel/opslevel-jq-parser/v2024"
	"github.com/rocktavious/autopilot/v2023"
)

//go:embed testdata/testService.json
var testServiceJSON []byte

//go:embed testdata/testRegistration.json
var testRegistrationJSON []byte

var gotService = false

type TestingClientBuilder struct {
	client *common.OpslevelClient
}

func NewTestingClientBuilder() *TestingClientBuilder {
	return &TestingClientBuilder{
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

func (b *TestingClientBuilder) SetAssignPropertyHandler(fn func(input opslevel.PropertyInput) error) *TestingClientBuilder {
	b.client.AssignPropertyHandler = fn
	return b
}

func (b *TestingClientBuilder) SetAssignTagsHandler(fn func(service *opslevel.Service, tags map[string]string) error) *TestingClientBuilder {
	b.client.AssignTagsHandler = fn
	return b
}

func (b *TestingClientBuilder) SetCreateAliasHandler(fn func(input opslevel.AliasCreateInput) error) *TestingClientBuilder {
	b.client.CreateAliasHandler = fn
	return b
}

func (b *TestingClientBuilder) SetCreateServiceHandler(fn func(input opslevel.ServiceCreateInput) (*opslevel.Service, error)) *TestingClientBuilder {
	b.client.CreateServiceHandler = fn
	return b
}

func (b *TestingClientBuilder) SetCreateServiceRepositoryHandler(fn func(input opslevel.ServiceRepositoryCreateInput) error) *TestingClientBuilder {
	b.client.CreateServiceRepositoryHandler = fn
	return b
}

func (b *TestingClientBuilder) SetCreateTagHandler(fn func(input opslevel.TagCreateInput) error) *TestingClientBuilder {
	b.client.CreateTagHandler = fn
	return b
}

func (b *TestingClientBuilder) SetCreateToolHandler(fn func(input opslevel.ToolCreateInput) error) *TestingClientBuilder {
	b.client.CreateToolHandler = fn
	return b
}

func (b *TestingClientBuilder) SetGetRepositoryWithAliasHandler(fn func(alias string) (*opslevel.Repository, error)) *TestingClientBuilder {
	b.client.GetRepositoryWithAliasHandler = fn
	return b
}

func (b *TestingClientBuilder) SetGetServiceHandler(fn func(alias string) (*opslevel.Service, error)) *TestingClientBuilder {
	b.client.GetServiceHandler = fn
	return b
}

func (b *TestingClientBuilder) SetUpdateServiceHandler(fn func(input opslevel.ServiceUpdateInput) (*opslevel.Service, error)) *TestingClientBuilder {
	b.client.UpdateServiceHandler = fn
	return b
}

func (b *TestingClientBuilder) SetUpdateServiceRepositoryHandler(fn func(input opslevel.ServiceRepositoryUpdateInput) error) *TestingClientBuilder {
	b.client.UpdateServiceRepositoryHandler = fn
	return b
}

func (b *TestingClientBuilder) GetClient() *common.OpslevelClient {
	return b.client
}

func TestReconcilerReconcile(t *testing.T) {
	// Arrange
	type TestCase struct {
		assert                 func(t *testing.T, err error)
		client                 *common.OpslevelClient
		registration           opslevel_jq_parser.ServiceRegistration
		disableServiceCreation bool
	}
	var testService opslevel.Service
	err := json.Unmarshal(testServiceJSON, &testService)
	if err != nil {
		panic(err)
	}
	var testRegistration opslevel_jq_parser.ServiceRegistration
	err = json.Unmarshal(testRegistrationJSON, &testRegistration)
	if err != nil {
		panic(err)
	}
	cases := map[string]TestCase{
		"Missing Aliases Should Error": {
			registration: opslevel_jq_parser.ServiceRegistration{
				Name:    "test",
				Aliases: []string{},
			},
			client: NewTestingClientBuilder().GetClient(),
			assert: func(t *testing.T, err error) {
				autopilot.Equals(t, "[test] found 0 aliases from kubernetes data", err.Error())
			},
		},
		"Matching Alias Should Call Service Update": {
			registration: testRegistration,
			client: NewTestingClientBuilder().SetGetServiceHandler(func(alias string) (*opslevel.Service, error) {
				// TODO: this global var is a hack to get a service to be returned only once
				if gotService == true {
					return nil, nil
				}
				gotService = true
				return &testService, nil
			}).SetUpdateServiceHandler(func(input opslevel.ServiceUpdateInput) (*opslevel.Service, error) {
				return nil, fmt.Errorf("done")
			}).GetClient(),
			assert: func(t *testing.T, err error) {
				autopilot.Equals(t, "done", err.Error())
			},
		},
		"Multiple Matching Aliases Should Halt": {
			registration: testRegistration,
			client: NewTestingClientBuilder().SetGetServiceHandler(func(alias string) (*opslevel.Service, error) {
				return &testService, nil
			}).GetClient(),
			assert: func(t *testing.T, err error) {
				autopilot.Equals(t, "[test] found multiple services with aliases = [[test1 test2 test3]]. cannot know which service to target for update ... skipping reconciliation", err.Error())
			},
		},
		"API Error On Get Service Should Halt": {
			registration: opslevel_jq_parser.ServiceRegistration{
				Name:    "test",
				Aliases: []string{"test"},
			},
			client: NewTestingClientBuilder().SetGetServiceHandler(func(alias string) (*opslevel.Service, error) {
				return nil, fmt.Errorf("api error")
			}).GetClient(),
			assert: func(t *testing.T, err error) {
				autopilot.Equals(t, "[test] api error during service lookup by alias.  unable to guarantee service was found or not ... skipping reconciliation", err.Error())
			},
		},
		// TODO: same on update service
		"API Error On Create Service Should Halt": {
			registration: opslevel_jq_parser.ServiceRegistration{
				Name:    "test",
				Aliases: []string{"test"},
			},
			client: NewTestingClientBuilder().SetGetServiceHandler(func(alias string) (*opslevel.Service, error) {
				return nil, nil
			}).SetCreateServiceHandler(func(input opslevel.ServiceCreateInput) (*opslevel.Service, error) {
				return nil, fmt.Errorf("api error")
			}).GetClient(),
			assert: func(t *testing.T, err error) {
				autopilot.Equals(t, "api error", err.Error())
			},
		},
		"Happy Path": {
			registration: testRegistration,
			client: NewTestingClientBuilder().SetGetServiceHandler(func(alias string) (*opslevel.Service, error) {
				return nil, nil
			}).SetCreateServiceHandler(func(input opslevel.ServiceCreateInput) (*opslevel.Service, error) {
				return &testService, nil
			}).SetCreateAliasHandler(func(input opslevel.AliasCreateInput) error {
				return nil
			}).SetAssignTagsHandler(func(service *opslevel.Service, tags map[string]string) error {
				return nil
			}).SetCreateTagHandler(func(input opslevel.TagCreateInput) error {
				return nil
			}).SetCreateToolHandler(func(input opslevel.ToolCreateInput) error {
				return nil
			}).SetCreateServiceRepositoryHandler(func(input opslevel.ServiceRepositoryCreateInput) error {
				return nil
			}).SetAssignPropertyHandler(func(input opslevel.PropertyInput) error {
				return nil
			}).GetClient(),
			assert: func(t *testing.T, err error) {
				autopilot.Ok(t, err)
			},
		},
		"Happy Path Do Not Create Services": {
			registration: testRegistration,
			client: NewTestingClientBuilder().SetGetServiceHandler(func(alias string) (*opslevel.Service, error) {
				return nil, nil
			}).GetClient(),
			assert: func(t *testing.T, err error) {
				autopilot.Ok(t, err)
			},
			disableServiceCreation: true,
		},
	}
	// Act
	for k, tc := range cases {
		t.Run(k, func(t *testing.T) {
			if k == "Single Matching Alias Should Call Service Update" {
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
	props := map[string]string{
		"prop_bool":         "true",
		"prop_empty_object": "{}",
		"prop_empty_string": "",
		"prop_object":       `{"message":"hello world","condition":true}`,
		"prop_string":       "hello world",
	}
	registration := opslevel_jq_parser.ServiceRegistration{
		Aliases:    []string{"a_test_service_with_properties"},
		Properties: props,
	}
	service := opslevel.Service{
		ServiceId: opslevel.ServiceId{
			Id: opslevel.ID("Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS85NzAyMg"),
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
	for _, x := range results {
		def := *x.Definition.Alias
		expId, gotId := service.ServiceId.Id, *x.Owner.Id
		autopilot.Assert(t, gotId == expId, fmt.Sprintf("[%s] unexpected owner ID '%s' - does not match service ID '%s'", def, gotId, expId))

		if val, ok := props[def]; ok {
			value, err := opslevel.NewJSONInput(val)
			if err != nil {
				log.Error().Err(err).Msgf("[%s] Failed parsing property: '%s'", service.Name, def)
				continue
			}
			expVal := string(*value)
			gotVal := string(x.Value)
			autopilot.Assert(t, gotVal == expVal, fmt.Sprintf("[%s] expected value for to be: '%s' got: '%s'", def, expVal, gotVal))
		} else {
			autopilot.Ok(t, fmt.Errorf("unexpected property definition alias: '%s'", def))
		}
	}
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
