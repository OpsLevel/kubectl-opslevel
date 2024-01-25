package autopilot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

type TestRequest struct {
	Request  GraphqlQuery
	Response map[string]any
}

func (t *TestRequest) ResponseAsString() string {
	marshalledResponse, err := json.Marshal(t.Response)
	if err != nil {
		panic(err)
	}
	return string(marshalledResponse)
}

func NewTestRequest(request string, variables string, response string) TestRequest {
	testRequest := TestRequest{
		Request: GraphqlQuery{
			Query:     templatedString(request),
			Variables: templatedJson(variables),
		},
		Response: templatedJson(response),
	}
	return testRequest
}

func templatedJson(values string) map[string]any {
	parsedValues, err := Templater.Use(values)
	if err != nil {
		panic(err)
	}
	var valuesJSON map[string]any
	if err := json.Unmarshal([]byte(parsedValues), &valuesJSON); err != nil {
		panic(err)
	}

	return valuesJSON
}

func templatedString(contents string) string {
	data, err := Templater.ParseToBytes(contents)
	if err != nil {
		panic(fmt.Errorf("error parsing template: %s", err))
	}
	return strings.TrimSpace(data.String())
}

func TestRequestResponse(testRequest TestRequest) ResponseWriter {
	return JsonStringResponse(testRequest.ResponseAsString())
}

func TestRequestValidation(t *testing.T, request TestRequest) RequestValidation {
	return GraphQLQueryToJsonValidation(t, request.Request)
}

func RegisterPaginatedEndpoint(t *testing.T, url string, requests ...TestRequest) string {
	requestCount := 0
	Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		GraphQLQueryToJsonValidation(t, requests[requestCount].Request)(r)
		TestRequestResponse(requests[requestCount])(w)
		requestCount += 1
	})
	return Server.URL + url
}
