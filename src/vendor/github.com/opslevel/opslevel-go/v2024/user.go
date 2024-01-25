package opslevel

import (
	"fmt"
	"slices"
)

type UserId struct {
	Id    ID
	Email string
}

type User struct {
	UserId
	HTMLUrl string
	Name    string
	Role    UserRole
	// We cannot have this here because its breaks a TON of queries
	// Teams   *TeamIdConnection
}

type UserConnection struct {
	Nodes      []User
	PageInfo   PageInfo
	TotalCount int
}

func (u *User) ResourceId() ID {
	return u.Id
}

func (u *User) ResourceType() TaggableResource {
	return TaggableResourceUser
}

//#region Helpers

func NewUserIdentifier(value string) *UserIdentifierInput {
	if IsID(value) {
		return &UserIdentifierInput{
			Id: NewID(value),
		}
	}
	return &UserIdentifierInput{
		Email: RefOf(value),
	}
}

func (u *UserId) GetTags(client *Client, variables *PayloadVariables) (*TagConnection, error) {
	var q struct {
		Account struct {
			User struct {
				Tags TagConnection `graphql:"tags(after: $after, first: $first)"`
			} `graphql:"user(id: $user)"`
		}
	}
	if u.Id == "" {
		return nil, fmt.Errorf("Unable to get Tags, invalid User id: '%s'", u.Id)
	}
	if variables == nil {
		variables = client.InitialPageVariablesPointer()
	}
	(*variables)["user"] = u.Id
	if err := client.Query(&q, *variables, WithName("UserTagsList")); err != nil {
		return nil, err
	}
	for q.Account.User.Tags.PageInfo.HasNextPage {
		(*variables)["after"] = q.Account.User.Tags.PageInfo.End
		resp, err := u.GetTags(client, variables)
		if err != nil {
			return nil, err
		}
		// Add unique tags only
		for _, resp := range resp.Nodes {
			if !slices.Contains[[]Tag, Tag](q.Account.User.Tags.Nodes, resp) {
				q.Account.User.Tags.Nodes = append(q.Account.User.Tags.Nodes, resp)
			}
		}
		q.Account.User.Tags.PageInfo = resp.PageInfo
		q.Account.User.Tags.TotalCount += resp.TotalCount
	}

	return &q.Account.User.Tags, nil
}

func (u *User) Teams(client *Client, variables *PayloadVariables) (*TeamIdConnection, error) {
	var q struct {
		Account struct {
			User struct {
				Teams TeamIdConnection `graphql:"teams(after: $after, first: $first)"`
			} `graphql:"user(id: $user)"`
		}
	}
	if u.Id == "" {
		return nil, fmt.Errorf("unable to get teams, nil user id")
	}
	if variables == nil {
		variables = client.InitialPageVariablesPointer()
	}
	(*variables)["user"] = u.Id
	if err := client.Query(&q, *variables, WithName("UserTeamsList")); err != nil { // what goes in "" here and how is it derived?
		return nil, err
	}
	for q.Account.User.Teams.PageInfo.HasNextPage {
		(*variables)["after"] = q.Account.User.Teams.PageInfo.End
		conn, err := u.Teams(client, variables)
		if err != nil {
			return nil, err
		}
		q.Account.User.Teams.Nodes = append(q.Account.User.Teams.Nodes, conn.Nodes...)
		q.Account.User.Teams.PageInfo = conn.PageInfo
		q.Account.User.Teams.TotalCount += conn.TotalCount
	}
	return &q.Account.User.Teams, nil
}

//#endregion

//#region Create

func (client *Client) InviteUser(email string, input UserInput) (*User, error) {
	var m struct {
		Payload struct {
			User   User
			Errors []OpsLevelErrors
		} `graphql:"userInvite(email: $email input: $input)"`
	}
	v := PayloadVariables{
		"email": email,
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("UserInvite"))
	return &m.Payload.User, HandleErrors(err, m.Payload.Errors)
}

//#endregion

//#region Retrieve

func (client *Client) GetUser(value string) (*User, error) {
	var q struct {
		Account struct {
			User User `graphql:"user(input: $input)"`
		}
	}
	v := PayloadVariables{
		"input": *NewUserIdentifier(value),
	}
	err := client.Query(&q, v, WithName("UserGet"))
	return &q.Account.User, HandleErrors(err, nil)
}

func (client *Client) ListUsers(variables *PayloadVariables) (*UserConnection, error) {
	var q struct {
		Account struct {
			Users UserConnection `graphql:"users(after: $after, first: $first)"`
		}
	}
	if variables == nil {
		variables = client.InitialPageVariablesPointer()
	}

	if err := client.Query(&q, *variables, WithName("UserList")); err != nil {
		return nil, err
	}

	for q.Account.Users.PageInfo.HasNextPage {
		(*variables)["after"] = q.Account.Users.PageInfo.End
		resp, err := client.ListUsers(variables)
		if err != nil {
			return nil, err
		}
		q.Account.Users.Nodes = append(q.Account.Users.Nodes, resp.Nodes...)
		q.Account.Users.PageInfo = resp.PageInfo
		q.Account.Users.TotalCount += resp.TotalCount
	}
	return &q.Account.Users, nil
}

//#endregion

//#region Update

func (client *Client) UpdateUser(user string, input UserInput) (*User, error) {
	var m struct {
		Payload struct {
			User   User
			Errors []OpsLevelErrors
		} `graphql:"userUpdate(user: $user input: $input)"`
	}
	v := PayloadVariables{
		"user":  *NewUserIdentifier(user),
		"input": input,
	}
	err := client.Mutate(&m, v, WithName("UserUpdate"))
	return &m.Payload.User, HandleErrors(err, m.Payload.Errors)
}

//#endregion

//#region Delete

func (client *Client) DeleteUser(user string) error {
	var m struct {
		Payload struct {
			Errors []OpsLevelErrors
		} `graphql:"userDelete(user: $user)"`
	}
	v := PayloadVariables{
		"user": *NewUserIdentifier(user),
	}
	err := client.Mutate(&m, v, WithName("UserDelete"))
	return HandleErrors(err, m.Payload.Errors)
}

//#endregion
