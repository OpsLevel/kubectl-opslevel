package opslevel

import (
	"github.com/shurcooL/graphql"
)

type TagOwner string

const (
	TagOwnerService TagOwner = "Service"
	TagOwnerRepository  TagOwner = "Repository"
)

type TaggableResource string

const (
	TaggableResourceService TaggableResource = "Service" // Used to identify a Service.
	TaggableResourceRepository  TaggableResource = "Repository"  // Used to identify a Repository.
)

type Tag struct {
	Id graphql.ID `json:"id"`
	Owner TagOwner `json:"owner"`
	Key graphql.String `json:"key"`
	Value graphql.String `json:"value"`
}

type TagInput struct {
	Key string `json:"key"`
	Value string `json:"value"`
}

type TagAssignInput struct {
	Id graphql.ID `json:"id"`
	Alias string `json:"alias,omitempty"`
	Type TaggableResource `json:"type,omitempty"`
	Tags []TagInput `json:"tags"`
}

type TagCreateInput struct {
	Id graphql.ID `json:"id"`
	Alias string `json:"alias,omitempty"`
	Type TaggableResource `json:"type,omitempty"`
	Key string `json:"key"`
	Value string `json:"value"`
}

type TagUpdateInput struct {
	Id graphql.ID `json:"id"`
	Key string `json:"key"`
	Value string `json:"value"`
}

type TagDeleteInput struct {
	Id graphql.ID `json:"id"`
}

//#region Assign

func (client *Client) AssignTagsForAlias(alias string, tags map[string]string) ([]Tag, error) {
	input := TagAssignInput{
		Alias: alias,
		Tags: []TagInput{},
	}
	for key, value := range tags {
        input.Tags = append(input.Tags, TagInput{
			Key: key,
			Value: value,
		})
    }
	return client.AssignTags(input)
}

func (client *Client) AssignTagsForId(id graphql.ID, tags map[string]string) ([]Tag, error) {
	input := TagAssignInput{
		Id: graphql.ID(id),
		Tags: []TagInput{},
	}
	for key, value := range tags {
        input.Tags = append(input.Tags, TagInput{
			Key: key,
			Value: value,
		})
    }
	return client.AssignTags(input)
}

func (client *Client) AssignTags(input TagAssignInput) ([]Tag, error) {
	var m struct {
		Payload struct {
			Tags []Tag
			Errors []OpsLevelErrors
		} `graphql:"tagAssign(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	if err := client.Mutate(&m, v); err != nil {
		return nil, err
	}
	return m.Payload.Tags, FormatErrors(m.Payload.Errors)
}

//#endregion

//#region Create

func (client *Client) CreateTags(alias string, tags map[string]string) ([]Tag, error) {
	var output []Tag
	for key, value := range tags {
		input := TagCreateInput{
			Alias: alias,
			Key: key,
			Value: value,
		}
		newTag, err := client.CreateTag(input)
        if (err != nil) {
			// TODO: combind errors?
		} else {
			output = append(output, *newTag)
		}
    }
	return output, nil
}

func (client *Client) CreateTagsForId(id graphql.ID, tags map[string]string) ([]Tag, error) {
	var output []Tag
	for key, value := range tags {
		input := TagCreateInput{
			Id: id,
			Key: key,
			Value: value,
		}
		newTag, err := client.CreateTag(input)
        if (err != nil) {
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
			Tag Tag
			Errors []OpsLevelErrors
		} `graphql:"tagCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	if err := client.Mutate(&m, v); err != nil {
		return nil, err
	}
	return &m.Payload.Tag, FormatErrors(m.Payload.Errors)
}

//#endregion

//#region Update

func (client *Client) UpdateTag(input TagUpdateInput) (*Tag, error) {
	var m struct {
		Payload struct {
			Tag Tag
			Errors []OpsLevelErrors
		} `graphql:"tagUpdate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	if err := client.Mutate(&m, v); err != nil {
		return nil, err
	}
	return &m.Payload.Tag, FormatErrors(m.Payload.Errors)
}

//#endregion

//#region Delete

func (client *Client) DeleteTag(id graphql.ID) error {
	var m struct {
		Payload struct {
			Id graphql.ID `graphql:"deletedTagId"`
			Errors []OpsLevelErrors
		} `graphql:"tagDelete(input: $input)"`
	}
	v := PayloadVariables{
		"input": TagDeleteInput{Id: id},
	}
	if err := client.Mutate(&m, v); err != nil {
		return err
	}
	return FormatErrors(m.Payload.Errors)
}

//#endregion
