package opslevel

import (
	"fmt"
	"regexp"
)

type TagOwner string

const (
	TagOwnerService    TagOwner = "Service"
	TagOwnerRepository TagOwner = "Repository"
)

type TaggableResourceInterface interface {
	GetTags(*Client, *PayloadVariables) (*TagConnection, error)
	ResourceId() ID
	ResourceType() TaggableResource
}

var (
	TagKeyRegex    = regexp.MustCompile(`\A[a-z][0-9a-z_\.\/\\-]*\z`)
	TagKeyErrorMsg = "tag key name '%s' must start with a letter and be only lowercase alphanumerics, underscores, hyphens, periods, and slashes"
)

type Tag struct {
	Id    ID     `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TagConnection struct {
	Nodes      []Tag
	PageInfo   PageInfo
	TotalCount int
}

func (client *Client) GetTaggableResource(resourceType TaggableResource, identifier string) (TaggableResourceInterface, error) {
	var err error
	var taggableResource TaggableResourceInterface

	switch resourceType {
	case TaggableResourceService:
		if IsID(identifier) {
			taggableResource, err = client.GetService(ID(identifier))
		} else {
			taggableResource, err = client.GetServiceWithAlias(identifier)
		}
	case TaggableResourceRepository:
		if IsID(identifier) {
			taggableResource, err = client.GetRepository(ID(identifier))
		} else {
			taggableResource, err = client.GetRepositoryWithAlias(identifier)
		}
	case TaggableResourceTeam:
		if IsID(identifier) {
			taggableResource, err = client.GetTeam(ID(identifier))
		} else {
			taggableResource, err = client.GetTeamWithAlias(identifier)
		}
	case TaggableResourceDomain:
		taggableResource, err = client.GetDomain(identifier)
	case TaggableResourceInfrastructureresource:
		taggableResource, err = client.GetInfrastructure(identifier)
	case TaggableResourceSystem:
		taggableResource, err = client.GetSystem(identifier)
	case TaggableResourceUser:
		taggableResource, err = client.GetUser(identifier)
	default:
		return nil, fmt.Errorf("not a taggable resource type: %s" + string(resourceType))
	}

	if err != nil {
		return nil, err
	}
	return taggableResource, nil
}

func (t *TagConnection) GetTagById(tagId ID) (*Tag, error) {
	for _, tag := range t.Nodes {
		if tag.Id == tagId {
			return &tag, nil
		}
	}
	return nil, fmt.Errorf("Tag with ID '%s' not found.", tagId)
}

//#region Helpers

func ValidateTagKey(key string) error {
	if !TagKeyRegex.MatchString(key) {
		return fmt.Errorf(TagKeyErrorMsg, key)
	}
	return nil
}

//#region Assign

func (client *Client) AssignTags(identifier string, tags map[string]string) ([]Tag, error) {
	input := TagAssignInput{
		Tags: []TagInput{},
	}
	for key, value := range tags {
		if err := ValidateTagKey(key); err != nil {
			return nil, err
		}
		input.Tags = append(input.Tags, TagInput{
			Key:   key,
			Value: value,
		})
	}
	if IsID(identifier) {
		input.Id = NewID(identifier)
	} else {
		input.Alias = &identifier
	}
	return client.AssignTag(input)
}

func (client *Client) AssignTag(input TagAssignInput) ([]Tag, error) {
	var m struct {
		Payload struct {
			Tags   []Tag
			Errors []OpsLevelErrors
		} `graphql:"tagAssign(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("TagAssign"))
	return m.Payload.Tags, HandleErrors(err, m.Payload.Errors)
}

//#endregion

//#region Create

func (client *Client) CreateTags(identifier string, tags map[string]string) ([]Tag, error) {
	var output []Tag
	for key, value := range tags {
		if err := ValidateTagKey(key); err != nil {
			return nil, err
		}
		input := TagCreateInput{
			Key:   key,
			Value: value,
		}
		if IsID(identifier) {
			input.Id = NewID(identifier)
		} else {
			input.Alias = &identifier
		}
		newTag, err := client.CreateTag(input)
		if err != nil {
			// TODO: combind errors?
		} else {
			output = append(output, *newTag)
		}
	}
	return output, nil
}

func (client *Client) CreateTag(input TagCreateInput) (*Tag, error) {
	var m struct {
		Payload struct {
			Tag    Tag `json:"tag"`
			Errors []OpsLevelErrors
		} `graphql:"tagCreate(input: $input)"`
	}
	if err := ValidateTagKey(input.Key); err != nil {
		return nil, err
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("TagCreate"))
	return &m.Payload.Tag, HandleErrors(err, m.Payload.Errors)
}

//#endregion

//#region Update

func (client *Client) UpdateTag(input TagUpdateInput) (*Tag, error) {
	var m struct {
		Payload struct {
			Tag    Tag
			Errors []OpsLevelErrors
		} `graphql:"tagUpdate(input: $input)"`
	}
	if err := ValidateTagKey(*input.Key); err != nil {
		return nil, err
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("TagUpdate"))
	return &m.Payload.Tag, HandleErrors(err, m.Payload.Errors)
}

//#endregion

//#region Delete

func (client *Client) DeleteTag(id ID) error {
	var m struct {
		Payload struct {
			Errors []OpsLevelErrors
		} `graphql:"tagDelete(input: $input)"`
	}
	v := PayloadVariables{
		"input": TagDeleteInput{Id: id},
	}
	err := client.Mutate(&m, v, WithName("TagDelete"))
	return HandleErrors(err, m.Payload.Errors)
}

//#endregion
