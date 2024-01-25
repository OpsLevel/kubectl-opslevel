package autopilot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type GraphqlQuery struct {
	Query     string
	Variables map[string]interface{} `json:",omitempty"`
}

func ToJson(query GraphqlQuery) string {
	bytes, _ := json.Marshal(query)
	return string(bytes)
}

func Parse(r *http.Request) GraphqlQuery {
	output := GraphqlQuery{}
	defer r.Body.Close()
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return output
	}
	if err = json.Unmarshal(bytes, &output); err != nil {
		fmt.Printf("autopilot error: %s", err.Error())
	}
	return output
}

func GraphQLQueryToJsonValidation(t *testing.T, request GraphqlQuery) RequestValidation {
	return func(r *http.Request) {
		Equals(t, ToJson(request), ToJson(Parse(r)))
	}
}

func GraphQLQueryValidation(t *testing.T, exp string) RequestValidation {
	return func(r *http.Request) {
		Equals(t, exp, Parse(r).Query)
	}
}

func GraphQLQueryFixture(fixture string) GraphqlQuery {
	exp := GraphqlQuery{}
	if err := json.Unmarshal([]byte(TemplatedFixture(fixture)), &exp); err != nil {
		fmt.Printf("autopilot error: %s", err.Error())
	}
	return exp
}

func GraphQLQueryFixtureValidation(t *testing.T, fixture string) RequestValidation {
	return func(r *http.Request) {
		Equals(t, ToJson(GraphQLQueryFixture(fixture)), ToJson(Parse(r)))
	}
}
