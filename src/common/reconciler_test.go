package common_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rs/zerolog/log"

	"golang.org/x/exp/maps"

	"github.com/opslevel/kubectl-opslevel/common"
	"github.com/opslevel/opslevel-go/v2024"
	opslevel_jq_parser "github.com/opslevel/opslevel-jq-parser/v2024"
	"github.com/rocktavious/autopilot/v2023"
)

// TODO: use builder pattern for client so that every function that isn't explicitly stubbed will panic

func TestReconcilerReconcile(t *testing.T) {
	// Arrange
	type TestCase struct {
		registration opslevel_jq_parser.ServiceRegistration
		reconciler   *common.ServiceReconciler
		assert       func(t *testing.T, err error)
	}
	// TODO: update testService to have system (parent)
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
	// TODO: update testRegistration to have system
	testRegistration := opslevel_jq_parser.ServiceRegistration{
		Aliases:      []string{"test1", "test2", "test3"},
		Description:  "test",
		Framework:    "test",
		Language:     "test",
		Lifecycle:    "test",
		Name:         "test",
		Owner:        "test",
		Product:      "test",
		Repositories: []opslevel.ServiceRepositoryCreateInput{{Service: *opslevel.NewIdentifier(""), Repository: *opslevel.NewIdentifier(""), DisplayName: opslevel.RefOf(""), BaseDirectory: opslevel.RefOf("")}},
		TagAssigns:   []opslevel.TagInput{{Key: "foo", Value: "bar"}, {Key: "hello", Value: "world"}},
		TagCreates:   []opslevel.TagInput{{Key: "env", Value: "test"}},
		Tier:         "test",
		Tools:        []opslevel.ToolCreateInput{{Category: opslevel.ToolCategoryCode, DisplayName: "test", Url: "https://example.com", Environment: opslevel.RefOf("test")}},
	}
	testRegistrationChangesAliasesOnly := opslevel_jq_parser.ServiceRegistration{
		Aliases: []string{"aliases_has_update", "updates_only_aliases"},
	}
	testRegistrationChangesDescriptionOnly := opslevel_jq_parser.ServiceRegistration{
		Aliases:     []string{"test"},
		Description: "description has update",
	}
	testRegistrationChangesEveryField := opslevel_jq_parser.ServiceRegistration{
		Aliases:     []string{"test1", "test2", "test3"},
		Description: "changed_description",
		Framework:   "changed_framework",
		Language:    "changed_language",
		Lifecycle:   "changed_lifecycle",
		Name:        "changed_name",
		Owner:       "changed_owner",
		Product:     "changed_product",
		System:      "changed_system",
		Tier:        "changed_tier",
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
					// TODO: this is nt a "nil service", this is an empty service - this is confusing because the service can be returned as a pointer...
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
		// TODO: need test where service creation is disabled and the service already exists
		"Happy Path Do Not Create Services": {
			registration: testRegistration,
			reconciler: common.NewServiceReconciler(&common.OpslevelClient{
				// TODO: if we are testing if service creation is disabled shouldn't this return nil, instead of empty service?
				GetServiceHandler: func(alias string) (*opslevel.Service, error) {
					return &opslevel.Service{}, nil // This returns a nil service as if the alias lookup didn't find anything
				},
				CreateServiceHandler: func(input opslevel.ServiceCreateInput) (*opslevel.Service, error) {
					panic("should not be called")
				},
				// TODO: should this panic? it is valid for a service to be updated if service creation is not enabled
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
		"Update Path - Has No Changes": {
			registration: testRegistration,
			reconciler: common.NewServiceReconciler(&common.OpslevelClient{
				GetServiceHandler: func(alias string) (*opslevel.Service, error) {
					return &testService, nil
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
			}, false),
			assert: func(t *testing.T, err error) {
				autopilot.Ok(t, err)
			},
		},
		"Update Path - Changes Only Aliases But No Fields": {
			registration: testRegistrationChangesAliasesOnly,
			reconciler: common.NewServiceReconciler(&common.OpslevelClient{
				GetServiceHandler: func(alias string) (*opslevel.Service, error) {
					return &testService, nil
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
			}, false),
			assert: func(t *testing.T, err error) {
				autopilot.Ok(t, err)
			},
		},
		"Update Path - Changes Only Description Field": {
			registration: testRegistrationChangesDescriptionOnly,
			reconciler: common.NewServiceReconciler(&common.OpslevelClient{
				GetServiceHandler: func(alias string) (*opslevel.Service, error) {
					return &testService, nil
				},
				CreateServiceHandler: func(input opslevel.ServiceCreateInput) (*opslevel.Service, error) {
					panic("should not be called")
				},
				UpdateServiceHandler: func(input opslevel.ServiceUpdateInput) (*opslevel.Service, error) {
					expectedInput := opslevel.ServiceUpdateInput{
						Id:          input.Id,
						Description: opslevel.RefOf("description has update"),
					}
					if diff := cmp.Diff(expectedInput, input); diff != "" {
						log.Panic().Msgf("expected different update input\nexp: %s\ngot: %s\n", toJSON(expectedInput), toJSON(input))
					}
					// TODO: this should update the service description and then return, but that would require adding dozens of lines of code. need a builder for services.
					return &testService, nil
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
		"Update Path - Changes Every Field": {
			registration: testRegistrationChangesEveryField,
			reconciler: common.NewServiceReconciler(&common.OpslevelClient{
				GetServiceHandler: func(alias string) (*opslevel.Service, error) {
					return &testService, nil
				},
				CreateServiceHandler: func(input opslevel.ServiceCreateInput) (*opslevel.Service, error) {
					panic("should not be called")
				},
				UpdateServiceHandler: func(input opslevel.ServiceUpdateInput) (*opslevel.Service, error) {
					// TODO: since opslevel-go client cannot be mocked, exclude lifecyce/owner/tier expectation for now...
					expectedInput := opslevel.ServiceUpdateInput{
						Id:          input.Id,
						Description: opslevel.RefOf("changed_description"),
						Framework:   opslevel.RefOf("changed_framework"),
						Language:    opslevel.RefOf("changed_language"),
						// LifecycleAlias: opslevel.RefOf("changed_lifecycle"),
						Name: opslevel.RefOf("changed_name"),
						// OwnerInput: opslevel.NewIdentifier("changed_owner"),
						Parent:  opslevel.NewIdentifier("changed_system"),
						Product: opslevel.RefOf("changed_product"),
						// TierAlias:  opslevel.RefOf("changed_tier"),
					}
					if diff := cmp.Diff(expectedInput, input); diff != "" {
						log.Panic().Msgf("expected different update input\nexp: %s\ngot: %s\n", toJSON(expectedInput), toJSON(input))
					}
					// TODO: this should update the service description and then return, but that would require adding dozens of lines of code. need a builder for services.
					return &testService, nil
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
	}
	// Act
	autopilot.RunTableTests(t, cases, func(t *testing.T, test TestCase) {
		// Assert
		test.assert(t, test.reconciler.Reconcile(test.registration))
	})
}

func Test_Reconciler_RepoNotInOpsLevel(t *testing.T) {
	testService := opslevel.Service{
		ServiceId: opslevel.ServiceId{Id: "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS8xNzg5Nw", Aliases: []string{"test"}},
		Name:      "Test Service",
	}
	testRegistration := opslevel_jq_parser.ServiceRegistration{
		Aliases: []string{"test"},
		Repositories: []opslevel.ServiceRepositoryCreateInput{
			{
				BaseDirectory: opslevel.RefOf(""),
				DisplayName:   opslevel.RefOf("Repo not in OpsLevel"),
				Repository:    *opslevel.NewIdentifier("repo_not_in_opslevel"),
				Service:       *opslevel.NewIdentifier(""),
			},
		},
	}
	calledGetRepositoryWithAliasHandler := false
	reconciler := common.NewServiceReconciler(&common.OpslevelClient{
		AssignPropertyHandler: func(input opslevel.PropertyInput) error {
			panic("should not be called")
		},
		AssignTagsHandler: func(service *opslevel.Service, tags map[string]string) error {
			panic("should not be called")
		},
		CreateAliasHandler: func(input opslevel.AliasCreateInput) error {
			panic("should not be called")
		},
		CreateServiceHandler: func(input opslevel.ServiceCreateInput) (*opslevel.Service, error) {
			panic("should not be called")
		},
		CreateServiceRepositoryHandler: func(input opslevel.ServiceRepositoryCreateInput) error {
			panic("should not be called")
		},
		CreateTagHandler: func(input opslevel.TagCreateInput) error {
			panic("should not be called")
		},
		CreateToolHandler: func(tool opslevel.ToolCreateInput) error {
			panic("should not be called")
		},
		GetRepositoryWithAliasHandler: func(alias string) (*opslevel.Repository, error) {
			calledGetRepositoryWithAliasHandler = true
			return nil, fmt.Errorf("[test] error fetching repository")
		},
		GetServiceHandler: func(alias string) (*opslevel.Service, error) {
			return &testService, nil
		},
		UpdateServiceHandler: func(input opslevel.ServiceUpdateInput) (*opslevel.Service, error) {
			panic("should not be called")
		},
		UpdateServiceRepositoryHandler: func(input opslevel.ServiceRepositoryUpdateInput) error {
			panic("should not be called")
		},
	}, true)
	reconcilerError := reconciler.Reconcile(testRegistration)

	autopilot.Ok(t, reconcilerError)
	autopilot.Assert(t, calledGetRepositoryWithAliasHandler, "expected call to GetRepositoryWithAliasHandler")
}

func Test_Reconciler_RepoIsAttached(t *testing.T) {
	testService := opslevel.Service{
		ServiceId: opslevel.ServiceId{Id: "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS8xNzg5Nw", Aliases: []string{"test"}},
		Name:      "Test Service",
	}
	testRegistration := opslevel_jq_parser.ServiceRegistration{
		Aliases: []string{"test"},
		Repositories: []opslevel.ServiceRepositoryCreateInput{
			{
				BaseDirectory: opslevel.RefOf(""),
				DisplayName:   opslevel.RefOf("Repo is attached"),
				Repository:    *opslevel.NewIdentifier("repo_is_attached"),
				Service:       *opslevel.NewIdentifier("test"),
			},
		},
	}
	calledGetRepositoryWithAliasHandler := false
	reconciler := common.NewServiceReconciler(&common.OpslevelClient{
		AssignPropertyHandler: func(input opslevel.PropertyInput) error {
			panic("should not be called")
		},
		AssignTagsHandler: func(service *opslevel.Service, tags map[string]string) error {
			panic("should not be called")
		},
		CreateAliasHandler: func(input opslevel.AliasCreateInput) error {
			panic("should not be called")
		},
		CreateServiceHandler: func(input opslevel.ServiceCreateInput) (*opslevel.Service, error) {
			panic("should not be called")
		},
		CreateServiceRepositoryHandler: func(input opslevel.ServiceRepositoryCreateInput) error {
			panic("should not be called")
		},
		CreateTagHandler: func(input opslevel.TagCreateInput) error {
			panic("should not be called")
		},
		CreateToolHandler: func(tool opslevel.ToolCreateInput) error {
			panic("should not be called")
		},
		GetRepositoryWithAliasHandler: func(alias string) (*opslevel.Repository, error) {
			calledGetRepositoryWithAliasHandler = true
			return &opslevel.Repository{
				Services: &opslevel.RepositoryServiceConnection{Edges: []opslevel.RepositoryServiceEdge{{
					Node: testService.ServiceId,
					ServiceRepositories: []opslevel.ServiceRepository{{
						BaseDirectory: "",
						DisplayName:   "Repo is attached",
						Id:            "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZVJlcG9zaXRvcnkvMTAwNDc",
						Repository:    opslevel.RepositoryId{DefaultAlias: "repo_is_attached"},
						Service:       testService.ServiceId,
					}},
				}}},
			}, nil
		},
		GetServiceHandler: func(alias string) (*opslevel.Service, error) {
			return &testService, nil
		},
		UpdateServiceHandler: func(input opslevel.ServiceUpdateInput) (*opslevel.Service, error) {
			panic("should not be called")
		},
		UpdateServiceRepositoryHandler: func(input opslevel.ServiceRepositoryUpdateInput) error {
			panic("should not be called")
		},
	}, true)
	reconcilerError := reconciler.Reconcile(testRegistration)

	autopilot.Ok(t, reconcilerError)
	autopilot.Assert(t, calledGetRepositoryWithAliasHandler, "expected call to GetRepositoryWithAliasHandler")
}

func Test_Reconciler_RepoNeedsCreate(t *testing.T) {
	testService := opslevel.Service{
		ServiceId: opslevel.ServiceId{Id: "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS8xNzg5Nw", Aliases: []string{"test"}},
		Name:      "Test Service",
	}
	testRegistration := opslevel_jq_parser.ServiceRegistration{
		Aliases: []string{"test"},
		Repositories: []opslevel.ServiceRepositoryCreateInput{
			{
				BaseDirectory: opslevel.RefOf(""),
				DisplayName:   opslevel.RefOf("Repo needs create"),
				Repository:    *opslevel.NewIdentifier("repo_needs_create"),
				Service:       *opslevel.NewIdentifier("test"),
			},
		},
	}
	calledGetRepositoryWithAliasHandler := false
	calledCreateServiceRepositoryHandler := false
	reconciler := common.NewServiceReconciler(&common.OpslevelClient{
		AssignPropertyHandler: func(input opslevel.PropertyInput) error {
			panic("should not be called")
		},
		AssignTagsHandler: func(service *opslevel.Service, tags map[string]string) error {
			panic("should not be called")
		},
		CreateAliasHandler: func(input opslevel.AliasCreateInput) error {
			panic("should not be called")
		},
		CreateServiceHandler: func(input opslevel.ServiceCreateInput) (*opslevel.Service, error) {
			panic("should not be called")
		},
		CreateServiceRepositoryHandler: func(input opslevel.ServiceRepositoryCreateInput) error {
			calledCreateServiceRepositoryHandler = true
			autopilot.Equals(t, *testRegistration.Repositories[0].BaseDirectory, *input.BaseDirectory)
			autopilot.Equals(t, *testRegistration.Repositories[0].DisplayName, *input.DisplayName)
			autopilot.Equals(t, "repo_needs_create", *input.Repository.Alias)
			autopilot.Equals(t, "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS8xNzg5Nw", string(*input.Service.Id))
			return nil
		},
		CreateTagHandler: func(input opslevel.TagCreateInput) error {
			panic("should not be called")
		},
		CreateToolHandler: func(tool opslevel.ToolCreateInput) error {
			panic("should not be called")
		},
		GetRepositoryWithAliasHandler: func(alias string) (*opslevel.Repository, error) {
			calledGetRepositoryWithAliasHandler = true
			// TODO: because we don't have a repo.GetService() handler, this stub is super dependent on the current logic of opslevel-go
			return &opslevel.Repository{
				Services: &opslevel.RepositoryServiceConnection{Edges: make([]opslevel.RepositoryServiceEdge, 0)},
			}, nil
		},
		GetServiceHandler: func(alias string) (*opslevel.Service, error) {
			return &testService, nil
		},
		UpdateServiceHandler: func(input opslevel.ServiceUpdateInput) (*opslevel.Service, error) {
			panic("should not be called")
		},
		UpdateServiceRepositoryHandler: func(input opslevel.ServiceRepositoryUpdateInput) error {
			panic("should not be called")
		},
	}, true)
	reconcilerError := reconciler.Reconcile(testRegistration)

	autopilot.Ok(t, reconcilerError)
	autopilot.Assert(t, calledGetRepositoryWithAliasHandler, "expected call to GetRepositoryWithAliasHandler")
	autopilot.Assert(t, calledCreateServiceRepositoryHandler, "expected call to CreateServiceRepositoryHandler")
}

func Test_Reconciler_RepoNeedsUpdate(t *testing.T) {
	testService := opslevel.Service{
		ServiceId: opslevel.ServiceId{Id: "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS8xNzg5Nw", Aliases: []string{"test"}},
		Name:      "Test Service",
	}
	testRegistration := opslevel_jq_parser.ServiceRegistration{
		Aliases: []string{"test"},
		Repositories: []opslevel.ServiceRepositoryCreateInput{
			{
				BaseDirectory: opslevel.RefOf(""),
				DisplayName:   opslevel.RefOf("Repo needs update"),
				Repository:    *opslevel.NewIdentifier("repo_needs_update"),
				Service:       *opslevel.NewIdentifier("test"),
			},
		},
	}
	calledGetRepositoryWithAliasHandler := false
	calledUpdateServiceRepositoryHandler := false
	reconciler := common.NewServiceReconciler(&common.OpslevelClient{
		AssignPropertyHandler: func(input opslevel.PropertyInput) error {
			panic("should not be called")
		},
		AssignTagsHandler: func(service *opslevel.Service, tags map[string]string) error {
			panic("should not be called")
		},
		CreateAliasHandler: func(input opslevel.AliasCreateInput) error {
			panic("should not be called")
		},
		CreateServiceHandler: func(input opslevel.ServiceCreateInput) (*opslevel.Service, error) {
			panic("should not be called")
		},
		CreateServiceRepositoryHandler: func(input opslevel.ServiceRepositoryCreateInput) error {
			panic("should not be called")
		},
		CreateTagHandler: func(input opslevel.TagCreateInput) error {
			panic("should not be called")
		},
		CreateToolHandler: func(tool opslevel.ToolCreateInput) error {
			panic("should not be called")
		},
		GetRepositoryWithAliasHandler: func(alias string) (*opslevel.Repository, error) {
			calledGetRepositoryWithAliasHandler = true
			return &opslevel.Repository{
				Services: &opslevel.RepositoryServiceConnection{Edges: []opslevel.RepositoryServiceEdge{
					{
						Node: testService.ServiceId,
						ServiceRepositories: []opslevel.ServiceRepository{
							{
								BaseDirectory: "",
								DisplayName:   "Hello World",
								Id:            "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZVJlcG9zaXRvcnkvMTAwNDc",
								Repository:    opslevel.RepositoryId{DefaultAlias: "repo_needs_update"},
								Service:       testService.ServiceId,
							},
						},
					},
				}},
			}, nil
		},
		GetServiceHandler: func(alias string) (*opslevel.Service, error) {
			return &testService, nil
		},
		UpdateServiceHandler: func(input opslevel.ServiceUpdateInput) (*opslevel.Service, error) {
			panic("should not be called")
		},
		UpdateServiceRepositoryHandler: func(input opslevel.ServiceRepositoryUpdateInput) error {
			calledUpdateServiceRepositoryHandler = true
			autopilot.Equals(t, *testRegistration.Repositories[0].BaseDirectory, *input.BaseDirectory)
			autopilot.Equals(t, *testRegistration.Repositories[0].DisplayName, *input.DisplayName)
			autopilot.Equals(t, "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZVJlcG9zaXRvcnkvMTAwNDc", string(input.Id))
			return nil
		},
	}, true)
	reconcilerError := reconciler.Reconcile(testRegistration)

	autopilot.Ok(t, reconcilerError)
	autopilot.Assert(t, calledGetRepositoryWithAliasHandler, "expected call to GetRepositoryWithAliasHandler")
	autopilot.Assert(t, calledUpdateServiceRepositoryHandler, "expected call to UpdateServiceRepositoryHandler")
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
	// testProperties is a map of strings to property inputs where the key is just the definition alias so that we can look it up easily in the expectation.
	testProperties := map[string]opslevel.PropertyInput{
		"prop_bool": {
			Definition: *opslevel.NewIdentifier("prop_bool"),
			Value:      opslevel.JsonString("true"),
		},
		"prop_empty_object": {
			Definition: *opslevel.NewIdentifier("prop_empty_object"),
			Value:      opslevel.JsonString("{}"),
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
	// for every property input send to the API by the reconciler, look up what it was originally (using the definition alias)
	// ensure that owner was set, definition did not change, value did not change
	for _, resultProperty := range results {
		autopilot.Equals(t, service.Id, *resultProperty.Owner.Id)
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

func toJSON[T any](object T) string {
	b, _ := json.Marshal(object)
	return string(b)
}
