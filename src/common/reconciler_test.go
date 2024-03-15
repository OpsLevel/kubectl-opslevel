package common_test

import (
	"fmt"
	"golang.org/x/exp/maps"
	"testing"

	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/opslevel/opslevel-go/v2024"
	opslevel_jq_parser "github.com/opslevel/opslevel-jq-parser/v2024"
	"github.com/rocktavious/autopilot/v2023"
)

func TestReconcilerReconcile(t *testing.T) {
	// Arrange
	type TestCase struct {
		registration opslevel_jq_parser.ServiceRegistration
		reconciler   *common.ServiceReconciler
		assert       func(t *testing.T, err error)
	}
	testService := opslevel.Service{
		ServiceId: opslevel.ServiceId{
			Id:      opslevel.ID("test"),
			Aliases: []string{"test"},
		},
		Name:        "test",
		Description: "test",
		Owner: opslevel.TeamId{
			Alias: "test",
		},
		Lifecycle: opslevel.Lifecycle{
			Alias: "test",
		},
		Tier: opslevel.Tier{
			Alias: "test",
		},
		Product:   "test",
		Language:  "test",
		Framework: "test",
		Tags: &opslevel.TagConnection{
			Nodes: []opslevel.Tag{
				{Key: "foo", Value: "bar"},
				{Key: "hello", Value: "world"},
				{Key: "env", Value: "test"},
			},
		},
		Tools: &opslevel.ToolConnection{
			Nodes: []opslevel.Tool{
				{Category: opslevel.ToolCategoryCode, DisplayName: "test", Url: "https://example.com", Environment: "test"},
			},
		},
	}
	testRegistration := opslevel_jq_parser.ServiceRegistration{
		Name:         "test",
		Description:  "test",
		Owner:        "test",
		Lifecycle:    "test",
		Tier:         "test",
		Product:      "test",
		Language:     "test",
		Framework:    "test",
		Aliases:      []string{"test1", "test2", "test3"},
		TagAssigns:   []opslevel.TagInput{{Key: "foo", Value: "bar"}, {Key: "hello", Value: "world"}},
		TagCreates:   []opslevel.TagInput{{Key: "env", Value: "test"}},
		Tools:        []opslevel.ToolCreateInput{{Category: opslevel.ToolCategoryCode, DisplayName: "test", Url: "https://example.com", Environment: opslevel.RefOf("test")}},
		Repositories: []opslevel.ServiceRepositoryCreateInput{{Service: *opslevel.NewIdentifier(""), Repository: *opslevel.NewIdentifier(""), DisplayName: opslevel.RefOf(""), BaseDirectory: opslevel.RefOf("")}},
	}
	cases := map[string]TestCase{
		"Missing Aliases Should Error": {
			registration: opslevel_jq_parser.ServiceRegistration{
				Name:    "test",
				Aliases: []string{},
			},
			reconciler: common.NewServiceReconciler(&common.OpslevelClient{}, false),
			assert: func(t *testing.T, err error) {
				autopilot.Equals(t, "[test] found 0 aliases from kubernetes data", err.Error())
			},
		},
		"Matching Alias Should Call Service Update": {
			registration: testRegistration,
			reconciler: common.NewServiceReconciler(&common.OpslevelClient{
				GetServiceHandler: func(alias string) (*opslevel.Service, error) {
					return &testService, nil
				},
				CreateServiceHandler: func(input opslevel.ServiceCreateInput) (*opslevel.Service, error) {
					panic("should not be called")
				},
				UpdateServiceHandler: func(input opslevel.ServiceUpdateInput) (*opslevel.Service, error) {
					return &testService, nil
				},
				GetRepositoryWithAliasHandler: func(alias string) (*opslevel.Repository, error) {
					return nil, fmt.Errorf("api error")
				},
			}, false),
			assert: func(t *testing.T, err error) {
				autopilot.Ok(t, err)
			},
		},
		"API Error On Get Service Should Halt": {
			registration: opslevel_jq_parser.ServiceRegistration{
				Name:    "test",
				Aliases: []string{"test"},
			},
			reconciler: common.NewServiceReconciler(&common.OpslevelClient{
				GetServiceHandler: func(alias string) (*opslevel.Service, error) {
					return nil, fmt.Errorf("api error")
				},
			}, false),
			assert: func(t *testing.T, err error) {
				autopilot.Equals(t, "[test] api error during service lookup by alias.  unable to guarantee service was found or not ... skipping reconciliation", err.Error())
			},
		},
		"API Error On Create Service Should Halt": {
			registration: opslevel_jq_parser.ServiceRegistration{
				Name:    "test",
				Aliases: []string{"test"},
			},
			reconciler: common.NewServiceReconciler(&common.OpslevelClient{
				GetServiceHandler: func(alias string) (*opslevel.Service, error) {
					return &testService, nil
				},
				CreateServiceHandler: func(input opslevel.ServiceCreateInput) (*opslevel.Service, error) {
					return nil, fmt.Errorf("api error")
				},
				CreateAliasHandler: func(input opslevel.AliasCreateInput) error {
					panic("should not be called")
				},
				AssignTagsHandler: func(service *opslevel.Service, tags map[string]string) error {
					panic("should not be called")
				},
				CreateTagHandler: func(input opslevel.TagCreateInput) error {
					panic("should not be called")
				},
			}, false),
			assert: func(t *testing.T, err error) {
				autopilot.Ok(t, err)
			},
		},
		"Happy Path": {
			registration: testRegistration,
			reconciler: common.NewServiceReconciler(&common.OpslevelClient{
				GetServiceHandler: func(alias string) (*opslevel.Service, error) {
					return &opslevel.Service{}, nil // This returns a nil service as if the alias lookup didn't find anything
				},
				CreateServiceHandler: func(input opslevel.ServiceCreateInput) (*opslevel.Service, error) {
					return &testService, nil
				},
				UpdateServiceHandler: func(input opslevel.ServiceUpdateInput) (*opslevel.Service, error) {
					panic("should not be called")
				},
				CreateAliasHandler: func(input opslevel.AliasCreateInput) error {
					return nil
				},
				AssignTagsHandler: func(service *opslevel.Service, tags map[string]string) error {
					return nil
				},
				CreateTagHandler: func(input opslevel.TagCreateInput) error {
					return nil
				},
				CreateToolHandler: func(tool opslevel.ToolCreateInput) error {
					return nil
				},
				GetRepositoryWithAliasHandler: func(alias string) (*opslevel.Repository, error) {
					return nil, fmt.Errorf("api error")
				},
				CreateServiceRepositoryHandler: func(input opslevel.ServiceRepositoryCreateInput) error {
					panic("should not be called")
				},
				UpdateServiceRepositoryHandler: func(input opslevel.ServiceRepositoryUpdateInput) error {
					panic("should not be called")
				},
			}, false),
			assert: func(t *testing.T, err error) {
				autopilot.Ok(t, err)
			},
		},
		"Happy Path Do Not Create Services": {
			registration: testRegistration,
			reconciler: common.NewServiceReconciler(&common.OpslevelClient{
				GetServiceHandler: func(alias string) (*opslevel.Service, error) {
					return &opslevel.Service{}, nil // This returns a nil service as if the alias lookup didn't find anything
				},
				CreateServiceHandler: func(input opslevel.ServiceCreateInput) (*opslevel.Service, error) {
					panic("should not be called")
				},
				UpdateServiceHandler: func(input opslevel.ServiceUpdateInput) (*opslevel.Service, error) {
					panic("should not be called")
				},
				CreateAliasHandler: func(input opslevel.AliasCreateInput) error {
					return nil
				},
				AssignTagsHandler: func(service *opslevel.Service, tags map[string]string) error {
					return nil
				},
				CreateTagHandler: func(input opslevel.TagCreateInput) error {
					return nil
				},
				CreateToolHandler: func(tool opslevel.ToolCreateInput) error {
					return nil
				},
				GetRepositoryWithAliasHandler: func(alias string) (*opslevel.Repository, error) {
					return nil, fmt.Errorf("api error")
				},
				CreateServiceRepositoryHandler: func(input opslevel.ServiceRepositoryCreateInput) error {
					panic("should not be called")
				},
				UpdateServiceRepositoryHandler: func(input opslevel.ServiceRepositoryUpdateInput) error {
					panic("should not be called")
				},
			}, true),
			assert: func(t *testing.T, err error) {
				autopilot.Ok(t, err)
			},
		},
	}
	// Act
	autopilot.RunTableTests(t, cases, func(t *testing.T, test TestCase) {
		// Assert
		test.assert(t, test.reconciler.Reconcile(test.registration))
	})
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

func Test_Reconciler_ServiceNeedsUpdate(t *testing.T) {
	type TestCase struct {
		input   opslevel.ServiceUpdateInput
		service opslevel.Service
		result  bool
	}
	// Arrange
	service1 := opslevel.Service{
		ServiceId: opslevel.ServiceId{
			Id: *opslevel.NewID("XXX"),
		},
		Name:        "Test",
		Description: "Hello World",
		Language:    "Python",
		Tier: opslevel.Tier{
			Alias: "tier_1",
		},
	}
	service2 := opslevel.Service{
		ServiceId: opslevel.ServiceId{
			Id: *opslevel.NewID("XXX"),
		},
		Name:        "Test",
		Description: "Hello World",
		Language:    "Python",
		Tier: opslevel.Tier{
			Alias: "tier_1",
		},
	}
	reconciler := common.NewServiceReconciler(&common.OpslevelClient{}, false)
	cases := map[string]TestCase{
		"Is True When Input Differs 1": {
			input: opslevel.ServiceUpdateInput{
				Name: opslevel.RefOf("Test1"),
			},
			service: service1,
			result:  true,
		},
		"Is True When Input Differs 2": {
			input: opslevel.ServiceUpdateInput{
				Name:      opslevel.RefOf("Test"),
				Language:  opslevel.RefOf("Python"),
				TierAlias: opslevel.RefOf("tier_2"),
			},
			service: service1,
			result:  true,
		},
		"Is True When Input Differs 3": {
			input: opslevel.ServiceUpdateInput{
				Name:       opslevel.RefOf("Test"),
				Language:   opslevel.RefOf("Python"),
				OwnerInput: opslevel.NewIdentifier("platform"),
			},
			service: service1,
			result:  true,
		},
		"Is False When Input Matches 1": {
			input: opslevel.ServiceUpdateInput{
				Id: opslevel.NewID("XXX"),
			},
			service: service2,
			result:  false,
		},
		"Is False When Input Matches 2": {
			input: opslevel.ServiceUpdateInput{
				Name: opslevel.RefOf("Test"),
			},
			service: service2,
			result:  false,
		},
		"Is False When Input Matches 3": {
			input: opslevel.ServiceUpdateInput{
				Name:      opslevel.RefOf("Test"),
				Language:  opslevel.RefOf("Python"),
				TierAlias: opslevel.RefOf("tier_1"),
			},
			service: service2,
			result:  false,
		},
	}
	// Act
	autopilot.RunTableTests(t, cases, func(t *testing.T, test TestCase) {
		// Assert
		autopilot.Equals(t, test.result, reconciler.ServiceNeedsUpdate(test.input, &test.service))
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
	// testProperties is a map of strings to property inputs where the key is just the definition alias so that we can look it up easily in the expectation.
	testProperties := map[string]opslevel.PropertyInput{
		"prop_bool": {
			Definition: *opslevel.NewIdentifier("prop_bool"),
			Value:      opslevel.JsonString("true"),
		},
		"prop_empty_object": {
			Definition: *opslevel.NewIdentifier("prop_empty_object"),
			Value:      opslevel.JsonString("{"),
		},
		"prop_object": {
			Definition: *opslevel.NewIdentifier("prop_object"),
			Value:      opslevel.JsonString("{\"message\":\"hello world\",\"condition\":true}"),
		},
		"prop_string": {
			Definition: *opslevel.NewIdentifier("prop_string"),
			Value:      opslevel.JsonString("hello world"),
		},
	}
	registration := opslevel_jq_parser.ServiceRegistration{
		Aliases:    []string{"a_test_service_with_properties"},
		Name:       "A test service with properties",
		Properties: maps.Values(testProperties),
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
	// for every property input processed by the reconciler, look up what it was before Reconcile()
	// ensure that service alias was set, definition did not change, value did not change
	for _, resultProperty := range results {
		autopilot.Equals(t, string(service.Id), *resultProperty.Owner.Alias)
		key := *resultProperty.Definition.Alias
		autopilot.Equals(t, *testProperties[key].Definition.Alias, *resultProperty.Definition.Alias)
		autopilot.Equals(t, string(testProperties[key].Value), string(resultProperty.Value))
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
