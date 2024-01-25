package opslevel

type ServiceDocumentSource struct {
	IntegrationId     `graphql:"... on ApiDocIntegration"`
	ServiceRepository `graphql:"... on ServiceRepository"`
}

type ServiceDocument struct {
	Id         ID                    `graphql:"id" json:"id"`
	HtmlURL    string                `graphql:"htmlUrl" json:"htmUrl,omitempty"`
	Source     ServiceDocumentSource `graphql:"source" json:"source"`
	Timestamps Timestamps            `graphql:"timestamps" json:"timestamps"`
}

type ServiceDocumentContent struct {
	ServiceDocument
	Content string `graphql:"content" json:"content,omitempty"`
}

func (c *Client) ServiceApiDocSettingsUpdate(service string, docPath string, docSource *ApiDocumentSourceEnum) (*Service, error) {
	var m struct {
		Payload struct {
			Service Service
			Errors  []OpsLevelErrors
		} `graphql:"serviceApiDocSettingsUpdate(service: $service, apiDocumentPath: $docPath, preferredApiDocumentSource: $docSource)"`
	}
	v := PayloadVariables{
		"service":   *NewIdentifier(service),
		"docPath":   NullString(),
		"docSource": docSource,
	}
	if docPath != "" {
		v["docPath"] = RefOf(docPath)
	}
	err := c.Mutate(&m, v, WithName("ServiceApiDocSettingsUpdate"))
	return &m.Payload.Service, HandleErrors(err, m.Payload.Errors)
}
