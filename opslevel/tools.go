package opslevel

import (
	"github.com/shurcooL/graphql"
)

type ToolCategory string

const (
	ToolCategoryAdmin                 ToolCategory = "admin"
	ToolCategoryCode                  ToolCategory = "code"
	ToolCategoryContinuousIntegration ToolCategory = "continuous_integration"
	ToolCategoryDeployment            ToolCategory = "deployment"
	ToolCategoryErrors                ToolCategory = "errors"
	ToolCategoryFeatureFlag           ToolCategory = "feature_flag"
	ToolCategoryHealthChecks          ToolCategory = "health_checks"
	ToolCategoryIncidents             ToolCategory = "incidents"
	ToolCategoryLogs                  ToolCategory = "logs"
	ToolCategoryMetrics               ToolCategory = "metrics"
	ToolCategoryOrchestrator          ToolCategory = "orchestrator"
	ToolCategoryRunbooks              ToolCategory = "runbooks"
	ToolCategoryStatusPage            ToolCategory = "status_page"
	ToolCategoryWiki                  ToolCategory = "wiki"
	ToolCategoryOther                 ToolCategory = "other"
)

type Tool struct {
	Category        ToolCategory
	CategoryAlias   graphql.String
	DisplayCategory graphql.String `json:",omitempty"`
	DisplayName     graphql.String
	Environment     graphql.String
	Id              graphql.ID
	Locked          graphql.Boolean
	PlainId         graphql.Int
	Service         Service
	Url             graphql.String
}

type ToolCreateInput struct {
	Category     ToolCategory `json:"category"`
	DisplayName  string       `json:"displayName"`
	Url          string       `json:"url"`
	Environment  string       `json:"environment"`
	ServiceId    graphql.ID   `json:"serviceId,omitempty"`
	ServiceAlias string       `json:"serviceAlias,omitempty"`
}

//#region Create

func (client *Client) CreateTool(input ToolCreateInput) (*Tool, error) {
	// TODO: validate - Category, DisplayName & Url are non nil - or throw err
	var m struct {
		Payload struct {
			Tool   Tool
			Errors []OpsLevelErrors
		} `graphql:"toolCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	if err := client.Mutate(&m, v); err != nil {
		return nil, err
	}
	return &m.Payload.Tool, FormatErrors(m.Payload.Errors)
}

//#endregion

//#region Retrieve

type ListToolQuery struct {
	Account struct {
		Tools struct {
			Nodes      []Tool
			PageInfo   PageInfo
			TotalCount graphql.Int
		} `graphql:"tools(after: $after, first: $first, service: $service)"`
	}
}

func (q *ListToolQuery) Query(client *Client, service graphql.ID) error {
	var subQ ListToolQuery
	v := PayloadVariables{
		"after":   q.Account.Tools.PageInfo.End,
		"first":   graphql.Int(100),
		"service": service,
	}
	if err := client.Query(&subQ, v); err != nil {
		return err
	}
	if subQ.Account.Tools.PageInfo.HasNextPage {
		subQ.Query(client, service)
	}
	for _, tool := range subQ.Account.Tools.Nodes {
		q.Account.Tools.Nodes = append(q.Account.Tools.Nodes, tool)
	}
	return nil
}

func (client *Client) ListTools(service graphql.ID) ([]Tool, error) {
	q := ListToolQuery{}
	if err := q.Query(client, service); err != nil {
		return []Tool{}, err
	}
	return q.Account.Tools.Nodes, nil
}

func (client *Client) GetToolCount(service graphql.ID) (int, error) {
	var q struct {
		Account struct {
			Tools struct {
				Nodes      []Tool
				PageInfo   PageInfo
				TotalCount graphql.Int
			} `graphql:"tools(service: $service)"`
		}
	}
	v := PayloadVariables{
		"service": service,
	}
	if err := client.Query(&q, v); err != nil {
		return 0, err
	}
	return int(q.Account.Tools.TotalCount), nil
}

//#endregion
