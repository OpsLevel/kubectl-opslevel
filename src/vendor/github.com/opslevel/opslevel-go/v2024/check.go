package opslevel

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/relvacode/iso8601"
)

type CheckOwner struct {
	Team TeamId `graphql:"... on Team"`
	// User User `graphql:"... on User"` // TODO: will this be public?
}

type Check struct {
	Category    Category     `graphql:"category"`
	Description string       `graphql:"description"`
	Enabled     bool         `graphql:"enabled"`
	EnableOn    iso8601.Time `graphql:"enableOn"`
	Filter      Filter       `graphql:"filter"`
	Id          ID           `graphql:"id"`
	Level       Level        `graphql:"level"`
	Name        string       `graphql:"name"`
	Notes       string       `graphql:"notes: rawNotes"`
	Owner       CheckOwner   `graphql:"owner"`
	Type        CheckType    `graphql:"type"`

	// TODO: resort these alphabetically - It will require fixing all the test fixtures
	AlertSourceUsageCheckFragment `graphql:"... on AlertSourceUsageCheck"`
	CustomEventCheckFragment      `graphql:"... on CustomEventCheck"`
	HasRecentDeployCheckFragment  `graphql:"... on HasRecentDeployCheck"`
	ManualCheckFragment           `graphql:"... on ManualCheck"`
	RepositoryFileCheckFragment   `graphql:"... on RepositoryFileCheck"`
	RepositoryGrepCheckFragment   `graphql:"... on RepositoryGrepCheck"`
	RepositorySearchCheckFragment `graphql:"... on RepositorySearchCheck"`
	ServiceOwnershipCheckFragment `graphql:"... on ServiceOwnershipCheck"`
	ServicePropertyCheckFragment  `graphql:"... on ServicePropertyCheck"`
	TagDefinedCheckFragment       `graphql:"... on TagDefinedCheck"`
	ToolUsageCheckFragment        `graphql:"... on ToolUsageCheck"`
	HasDocumentationCheckFragment `graphql:"... on HasDocumentationCheck"`
}

type CheckInputConstructor func() any

var CheckCreateConstructors = map[CheckType]CheckInputConstructor{
	CheckTypeAlertSourceUsage:    func() any { return &CheckAlertSourceUsageCreateInput{} },
	CheckTypeCustom:              func() any { return &CheckCreateInput{} },
	CheckTypeGeneric:             func() any { return &CheckCustomEventCreateInput{} },
	CheckTypeGitBranchProtection: func() any { return &CheckGitBranchProtectionCreateInput{} },
	CheckTypeHasDocumentation:    func() any { return &CheckHasDocumentationCreateInput{} },
	CheckTypeHasOwner:            func() any { return &CheckServiceOwnershipCreateInput{} },
	CheckTypeHasRecentDeploy:     func() any { return &CheckHasRecentDeployCreateInput{} },
	CheckTypeHasRepository:       func() any { return &CheckRepositoryIntegratedCreateInput{} },
	CheckTypeHasServiceConfig:    func() any { return &CheckServiceConfigurationCreateInput{} },
	CheckTypeManual:              func() any { return &CheckManualCreateInput{} },
	CheckTypePayload:             func() any { return &CheckCreateInput{} },
	CheckTypeRepoFile:            func() any { return &CheckRepositoryFileCreateInput{} },
	CheckTypeRepoGrep:            func() any { return &CheckRepositoryGrepCreateInput{} },
	CheckTypeRepoSearch:          func() any { return &CheckRepositorySearchCreateInput{} },
	CheckTypeServiceDependency:   func() any { return &CheckServiceDependencyCreateInput{} },
	CheckTypeServiceProperty:     func() any { return &CheckServicePropertyCreateInput{} },
	CheckTypeTagDefined:          func() any { return &CheckTagDefinedCreateInput{} },
	CheckTypeToolUsage:           func() any { return &CheckToolUsageCreateInput{} },
}

var CheckUpdateConstructors = map[CheckType]CheckInputConstructor{
	CheckTypeAlertSourceUsage:    func() any { return &CheckAlertSourceUsageUpdateInput{} },
	CheckTypeCustom:              func() any { return &CheckUpdateInput{} },
	CheckTypeGeneric:             func() any { return &CheckCustomEventUpdateInput{} },
	CheckTypeGitBranchProtection: func() any { return &CheckGitBranchProtectionUpdateInput{} },
	CheckTypeHasDocumentation:    func() any { return &CheckHasDocumentationUpdateInput{} },
	CheckTypeHasOwner:            func() any { return &CheckServiceOwnershipUpdateInput{} },
	CheckTypeHasRecentDeploy:     func() any { return &CheckHasRecentDeployUpdateInput{} },
	CheckTypeHasRepository:       func() any { return &CheckRepositoryIntegratedUpdateInput{} },
	CheckTypeHasServiceConfig:    func() any { return &CheckServiceConfigurationUpdateInput{} },
	CheckTypeManual:              func() any { return &CheckManualUpdateInput{} },
	CheckTypePayload:             func() any { return &CheckUpdateInput{} },
	CheckTypeRepoFile:            func() any { return &CheckRepositoryFileUpdateInput{} },
	CheckTypeRepoGrep:            func() any { return &CheckRepositoryGrepUpdateInput{} },
	CheckTypeRepoSearch:          func() any { return &CheckRepositorySearchUpdateInput{} },
	CheckTypeServiceDependency:   func() any { return &CheckServiceDependencyUpdateInput{} },
	CheckTypeServiceProperty:     func() any { return &CheckServicePropertyUpdateInput{} },
	CheckTypeTagDefined:          func() any { return &CheckTagDefinedUpdateInput{} },
	CheckTypeToolUsage:           func() any { return &CheckToolUsageUpdateInput{} },
}

func UnmarshalCheckCreateInput(checkType CheckType, data []byte) (any, error) {
	output := CheckCreateConstructors[checkType]()
	if err := json.Unmarshal(data, &output); err != nil {
		return nil, err
	}
	return output, nil
}

func UnmarshalCheckUpdateInput(checkType CheckType, data []byte) (any, error) {
	output := CheckUpdateConstructors[checkType]()
	if err := json.Unmarshal(data, &output); err != nil {
		return nil, err
	}
	return output, nil
}

type CheckConnection struct {
	Nodes      []Check
	PageInfo   PageInfo
	TotalCount int
}

type CheckCreateInputProvider interface {
	GetCheckCreateInput() *CheckCreateInput
}

type CheckCreateInput struct {
	Name     string        `json:"name" yaml:"name" mapstructure:"name"`
	Enabled  *bool         `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	EnableOn *iso8601.Time `json:"enableOn,omitempty" yaml:"enableOn,omitempty" mapstructure:"enabledOn,omitempty"`
	Category ID            `json:"categoryId" yaml:"categoryId" mapstructure:"categoryId"`
	Level    ID            `json:"levelId" yaml:"levelId" mapstructure:"levelId"`
	Owner    *ID           `json:"ownerId,omitempty" yaml:"ownerId,omitempty" mapstructure:"ownerId,omitempty"`
	Filter   *ID           `json:"filterId,omitempty" yaml:"filterId,omitempty" mapstructure:"filterId,omitempty"`
	Notes    *string       `json:"notes" yaml:"notes" default:"Notes on Example Check" mapstructure:"notes"`
}

func NewCheckCreateInputTypeOf[T any](checkCreateInput CheckCreateInput) *T {
	newCheck := new(T)
	if err := mapstructure.Decode(checkCreateInput, newCheck); err != nil {
		panic(err)
	}
	return newCheck
}

func (c *CheckCreateInput) GetCheckCreateInput() *CheckCreateInput {
	return c
}

type CheckUpdateInputProvider interface {
	GetCheckUpdateInput() *CheckUpdateInput
}

type CheckUpdateInput struct {
	Id       ID            `json:"id" mapstructure:"id"`
	Name     string        `json:"name,omitempty" mapstructure:"name,omitempty"`
	Enabled  *bool         `json:"enabled,omitempty" mapstructure:"enabled,omitempty"`
	EnableOn *iso8601.Time `json:"enableOn,omitempty" mapstructure:"enabledOn,omitempty"`
	Category ID            `json:"categoryId,omitempty" mapstructure:"categoryId,omitempty"`
	Level    ID            `json:"levelId,omitempty" mapstructure:"levelId,omitempty"`
	Owner    *ID           `json:"ownerId,omitempty" yaml:"ownerId,omitempty" mapstructure:"ownerId,omitempty"`
	Filter   *ID           `json:"filterId,omitempty" yaml:"filterId,omitempty" mapstructure:"filterId,omitempty"`
	Notes    *string       `json:"notes,omitempty" mapstructure:"notes,omitempty"`
}

func NewCheckUpdateInputTypeOf[T any](checkUpdateInput CheckUpdateInput) *T {
	newCheck := new(T)
	if err := mapstructure.Decode(checkUpdateInput, newCheck); err != nil {
		panic(err)
	}
	return newCheck
}

func (c *CheckUpdateInput) GetCheckUpdateInput() *CheckUpdateInput {
	return c
}

// Encompass CheckCreatePayload and CheckUpdatePayload into 1 struct
type CheckResponsePayload struct {
	Check  Check
	Errors []OpsLevelErrors
}

//#region Create

func (client *Client) CreateCheck(input any) (*Check, error) {
	switch v := input.(type) {
	case *CheckAlertSourceUsageCreateInput:
		return client.CreateCheckAlertSourceUsage(*v)
	case *CheckCustomEventCreateInput:
		return client.CreateCheckCustomEvent(*v)
	case *CheckGitBranchProtectionCreateInput:
		return client.CreateCheckGitBranchProtection(*v)
	case *CheckHasDocumentationCreateInput:
		return client.CreateCheckHasDocumentation(*v)
	case *CheckServiceOwnershipCreateInput:
		return client.CreateCheckServiceOwnership(*v)
	case *CheckHasRecentDeployCreateInput:
		return client.CreateCheckHasRecentDeploy(*v)
	case *CheckRepositoryIntegratedCreateInput:
		return client.CreateCheckRepositoryIntegrated(*v)
	case *CheckServiceConfigurationCreateInput:
		return client.CreateCheckServiceConfiguration(*v)
	case *CheckManualCreateInput:
		return client.CreateCheckManual(*v)
	case *CheckRepositoryFileCreateInput:
		return client.CreateCheckRepositoryFile(*v)
	case *CheckRepositoryGrepCreateInput:
		return client.CreateCheckRepositoryGrep(*v)
	case *CheckRepositorySearchCreateInput:
		return client.CreateCheckRepositorySearch(*v)
	case *CheckServiceDependencyCreateInput:
		return client.CreateCheckServiceDependency(*v)
	case *CheckServicePropertyCreateInput:
		return client.CreateCheckServiceProperty(*v)
	case *CheckTagDefinedCreateInput:
		return client.CreateCheckTagDefined(*v)
	case *CheckToolUsageCreateInput:
		return client.CreateCheckToolUsage(*v)
	}
	return nil, fmt.Errorf("unknown input type %T", input)
}

// See files check_*.go

//#endregion

//#region Retrieve

func (client *Client) GetCheck(id ID) (*Check, error) {
	var q struct {
		Account struct {
			Check Check `graphql:"check(id: $id)"`
		}
	}
	v := PayloadVariables{
		"id": id,
	}
	err := client.Query(&q, v, WithName("CheckGet"))
	if q.Account.Check.Id == "" {
		err = fmt.Errorf("check with ID '%s' not found", id)
	}
	return &q.Account.Check, HandleErrors(err, nil)
}

func (client *Client) ListChecks(variables *PayloadVariables) (*CheckConnection, error) {
	var q struct {
		Account struct {
			Rubric struct {
				Checks CheckConnection `graphql:"checks(after: $after, first: $first)"`
			}
		}
	}
	if variables == nil {
		variables = client.InitialPageVariablesPointer()
	}
	if err := client.Query(&q, *variables, WithName("CheckList")); err != nil {
		return nil, err
	}
	for q.Account.Rubric.Checks.PageInfo.HasNextPage {
		(*variables)["after"] = q.Account.Rubric.Checks.PageInfo.End
		resp, err := client.ListChecks(variables)
		if err != nil {
			return nil, err
		}
		q.Account.Rubric.Checks.Nodes = append(q.Account.Rubric.Checks.Nodes, resp.Nodes...)
		q.Account.Rubric.Checks.PageInfo = resp.PageInfo
		q.Account.Rubric.Checks.TotalCount += resp.TotalCount
	}
	return &q.Account.Rubric.Checks, nil
}

//#endregion

//#region Update

func (client *Client) UpdateCheck(input any) (*Check, error) {
	switch v := input.(type) {
	case *CheckAlertSourceUsageUpdateInput:
		return client.UpdateCheckAlertSourceUsage(*v)
	case *CheckCustomEventUpdateInput:
		return client.UpdateCheckCustomEvent(*v)
	case *CheckGitBranchProtectionUpdateInput:
		return client.UpdateCheckGitBranchProtection(*v)
	case *CheckHasDocumentationUpdateInput:
		return client.UpdateCheckHasDocumentation(*v)
	case *CheckServiceOwnershipUpdateInput:
		return client.UpdateCheckServiceOwnership(*v)
	case *CheckHasRecentDeployUpdateInput:
		return client.UpdateCheckHasRecentDeploy(*v)
	case *CheckRepositoryIntegratedUpdateInput:
		return client.UpdateCheckRepositoryIntegrated(*v)
	case *CheckServiceConfigurationUpdateInput:
		return client.UpdateCheckServiceConfiguration(*v)
	case *CheckManualUpdateInput:
		return client.UpdateCheckManual(*v)
	case *CheckRepositoryFileUpdateInput:
		return client.UpdateCheckRepositoryFile(*v)
	case *CheckRepositoryGrepUpdateInput:
		return client.UpdateCheckRepositoryGrep(*v)
	case *CheckRepositorySearchUpdateInput:
		return client.UpdateCheckRepositorySearch(*v)
	case *CheckServiceDependencyUpdateInput:
		return client.UpdateCheckServiceDependency(*v)
	case *CheckServicePropertyUpdateInput:
		return client.UpdateCheckServiceProperty(*v)
	case *CheckTagDefinedUpdateInput:
		return client.UpdateCheckTagDefined(*v)
	case *CheckToolUsageUpdateInput:
		return client.UpdateCheckToolUsage(*v)
	}
	return nil, fmt.Errorf("unknown input type %T", input)
}

//#endregion

//#region Delete

func (client *Client) DeleteCheck(id ID) error {
	var m struct {
		Payload struct {
			Errors []OpsLevelErrors
		} `graphql:"checkDelete(input: $input)"`
	}
	v := PayloadVariables{
		"input": CheckDeleteInput{Id: &id},
	}
	err := client.Mutate(&m, v, WithName("CheckDelete"))
	return HandleErrors(err, m.Payload.Errors)
}

//#endregion
