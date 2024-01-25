package opslevel

import "fmt"

type ServiceDependency struct {
	Id        ID        `graphql:"id"`
	Service   ServiceId `graphql:"sourceService"`
	DependsOn ServiceId `graphql:"destinationService"`
	Notes     string    `graphql:"notes"`
}

type ServiceDependenciesEdge struct {
	Id     ID         `graphql:"id"`
	Locked bool       `graphql:"locked"`
	Node   *ServiceId `graphql:"node"`
	Notes  string     `graphql:"notes"`
}

type ServiceDependentsEdge struct {
	Id     ID         `graphql:"id"`
	Locked bool       `graphql:"locked"`
	Node   *ServiceId `graphql:"node"`
	Notes  string     `graphql:"notes"`
}

type ServiceDependenciesConnection struct {
	Edges    []ServiceDependenciesEdge `graphql:"edges"`
	PageInfo PageInfo
}

type ServiceDependentsConnection struct {
	Edges    []ServiceDependentsEdge `graphql:"edges"`
	PageInfo PageInfo
}

//#region Create

func (client *Client) CreateServiceDependency(input ServiceDependencyCreateInput) (*ServiceDependency, error) {
	var m struct {
		Payload struct {
			ServiceDependency *ServiceDependency
			Errors            []OpsLevelErrors
		} `graphql:"serviceDependencyCreate(inputV2: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("ServiceDependencyCreate"))
	return m.Payload.ServiceDependency, HandleErrors(err, m.Payload.Errors)
}

//#endregion

//#region Retrieve

func (s *Service) GetDependencies(client *Client, variables *PayloadVariables) (*ServiceDependenciesConnection, error) {
	var q struct {
		Account struct {
			Service struct {
				Dependencies ServiceDependenciesConnection `graphql:"dependencies(after: $after, first: $first)"`
			} `graphql:"service(id: $service)"`
		}
	}
	if s.Id == "" {
		return nil, fmt.Errorf("Unable to get Dependencies, invalid service id: '%s'", s.Id)
	}
	if variables == nil {
		variables = client.InitialPageVariablesPointer()
	}
	(*variables)["service"] = s.Id
	if err := client.Query(&q, *variables, WithName("ServiceDependenciesList")); err != nil {
		return nil, err
	}
	if s.Dependencies == nil {
		s.Dependencies = &ServiceDependenciesConnection{}
	}
	s.Dependencies.Edges = append(s.Dependencies.Edges, q.Account.Service.Dependencies.Edges...)
	s.Dependencies.PageInfo = q.Account.Service.Dependencies.PageInfo
	for s.Dependencies.PageInfo.HasNextPage {
		(*variables)["after"] = s.Dependencies.PageInfo.End
		_, err := s.GetDependencies(client, variables)
		if err != nil {
			return nil, err
		}
	}
	return s.Dependencies, nil
}

func (s *Service) GetDependents(client *Client, variables *PayloadVariables) (*ServiceDependentsConnection, error) {
	var q struct {
		Account struct {
			Service struct {
				Dependents ServiceDependentsConnection `graphql:"dependents(after: $after, first: $first)"`
			} `graphql:"service(id: $service)"`
		}
	}
	if s.Id == "" {
		return nil, fmt.Errorf("Unable to get Dependents, invalid service id: '%s'", s.Id)
	}
	if variables == nil {
		variables = client.InitialPageVariablesPointer()
	}
	(*variables)["service"] = s.Id
	if err := client.Query(&q, *variables, WithName("ServiceDependentsList")); err != nil {
		return nil, err
	}
	if s.Dependents == nil {
		s.Dependents = &ServiceDependentsConnection{}
	}
	s.Dependents.Edges = append(s.Dependents.Edges, q.Account.Service.Dependents.Edges...)
	s.Dependents.PageInfo = q.Account.Service.Dependents.PageInfo
	for s.Dependents.PageInfo.HasNextPage {
		(*variables)["after"] = s.Dependents.PageInfo.End
		_, err := s.GetDependents(client, variables)
		if err != nil {
			return nil, err
		}
	}
	return s.Dependents, nil
}

//#endregion

//#region Delete

func (client *Client) DeleteServiceDependency(id ID) error {
	var m struct {
		Payload struct {
			Errors []OpsLevelErrors
		} `graphql:"serviceDependencyDelete(input: $input)"`
	}
	v := PayloadVariables{
		"input": DeleteInput{Id: id},
	}
	err := client.Mutate(&m, v, WithName("ServiceDependencyDelete"))
	return HandleErrors(err, m.Payload.Errors)
}

//#endregion
