package common

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/opslevel/opslevel-go"
	"github.com/rocktavious/autopilot"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Helper Functions
type MockedResponse interface {
	GetStatus() int
	GetResponse() []byte
}

type ObjectMockResponse struct {
	Status int
	Object interface{}
}

func (s ObjectMockResponse) GetStatus() int { return s.Status }
func (s ObjectMockResponse) GetResponse() []byte {
	data, err := json.Marshal(s.Object)
	// TODO: should this just log the problem and return "{}"
	if err != nil {
		panic(err)
	}
	return data
}

type FixtureMockResponse struct {
	Status int
	Path   string
}

func (s FixtureMockResponse) GetStatus() int { return s.Status }
func (s FixtureMockResponse) GetResponse() []byte {
	return []byte(autopilot.Fixture(fmt.Sprintf("%s.json", s.Path)))
}

type StringMockResponse struct {
	Status int
	Data   string
}

func (s StringMockResponse) GetStatus() int      { return s.Status }
func (s StringMockResponse) GetResponse() []byte { return []byte(s.Data) }

func AMockedClient(responses ...MockedResponse) (*opslevel.Client, *httptest.Server) {
	count := len(responses)
	current := 0
	mockedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if current < count {
			response := responses[current]
			w.WriteHeader(response.GetStatus())
			w.Write(response.GetResponse())
			current += 1
		} else {
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintf(w, "{}")
		}
	}))
	mockedClient := opslevel.NewClient("X", opslevel.SetURL(mockedServer.URL))
	return mockedClient, mockedServer
}

//////////////////////////

func TestMain(m *testing.M) {
	output := zerolog.ConsoleWriter{Out: os.Stderr}
	log.Logger = log.Output(output)
	flag.Parse()
	teardown := autopilot.Setup()
	defer teardown()
	os.Exit(m.Run())
}

func Test_ContainsAllTags_IsTrue_WhenAllTagsOverlap(t *testing.T) {
	// Arrange
	inputA := []opslevel.TagInput{
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
	}
	inputB := []opslevel.Tag{
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
	}
	// Act
	result := containsAllTags(inputA, inputB)
	// Assert
	autopilot.Equals(t, true, result)
}

func Test_ContainsAllTags_IsFalse_WhenTagInputHasMore(t *testing.T) {
	// Arrange
	inputA := []opslevel.TagInput{
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
	}
	inputB := []opslevel.Tag{
		{
			Key:   "foo",
			Value: "bar",
		},
		{
			Key:   "hello",
			Value: "world",
		},
	}
	// Act
	result := containsAllTags(inputA, inputB)
	// Assert
	autopilot.Equals(t, false, result)
}

func Test_ContainsAllTags_IsTrue_WhenServiceHasMore(t *testing.T) {
	// Arrange
	inputA := []opslevel.TagInput{
		{
			Key:   "foo",
			Value: "bar",
		},
	}
	inputB := []opslevel.Tag{
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
	}
	// Act
	result := containsAllTags(inputA, inputB)
	// Assert
	autopilot.Equals(t, true, result)
}

func Test_ServiceNeedsUpdate_IsTrue_WhenInputDiffers(t *testing.T) {
	// Arrange
	service := opslevel.Service{
		ServiceId: opslevel.ServiceId{
			Id: opslevel.NewID("XXX"),
		},
		Name:        "Test",
		Description: "Hello World",
		Language:    "Python",
		Tier: opslevel.Tier{
			Alias: "tier_1",
		},
	}
	input1 := opslevel.ServiceUpdateInput{
		Name: "Test1",
	}
	input2 := opslevel.ServiceUpdateInput{
		Name:     "Test",
		Language: "Python",
		Tier:     "tier_2",
	}
	input3 := opslevel.ServiceUpdateInput{
		Name:     "Test",
		Language: "Python",
		Owner:    "platform",
	}
	// Act
	result1 := serviceNeedsUpdate(input1, &service)
	result2 := serviceNeedsUpdate(input2, &service)
	result3 := serviceNeedsUpdate(input3, &service)
	// Assert
	autopilot.Equals(t, true, result1)
	autopilot.Equals(t, true, result2)
	autopilot.Equals(t, true, result3)
}

func Test_ServiceNeedsUpdate_IsFalse_WhenInputMatches(t *testing.T) {
	// Arrange
	service := opslevel.Service{
		ServiceId: opslevel.ServiceId{
			Id: opslevel.NewID("XXX"),
		},
		Name:        "Test",
		Description: "Hello World",
		Language:    "Python",
		Tier: opslevel.Tier{
			Alias: "tier_1",
		},
	}
	input1 := opslevel.ServiceUpdateInput{
		Id: "XXX",
	}
	input2 := opslevel.ServiceUpdateInput{
		Name: "Test",
	}
	input3 := opslevel.ServiceUpdateInput{
		Name:     "Test",
		Language: "Python",
		Tier:     "tier_1",
	}
	// Act
	result1 := serviceNeedsUpdate(input1, &service)
	result2 := serviceNeedsUpdate(input2, &service)
	result3 := serviceNeedsUpdate(input3, &service)
	// Assert
	autopilot.Equals(t, false, result1)
	autopilot.Equals(t, false, result2)
	autopilot.Equals(t, false, result3)
}

func Test_ValidateServiceAliases_WhenNoAliasesMatch(t *testing.T) {
	// Arrange
	mockedResponse := StringMockResponse{
		Status: http.StatusOK,
		Data:   "{}",
	}
	mockedClient, mockedServer := AMockedClient(mockedResponse, mockedResponse, mockedResponse)
	defer mockedServer.Close()
	registration := ServiceRegistration{
		Name: "Test",
		Aliases: []string{
			"Alias1",
			"Alias2",
			"Alias3",
		},
	}
	// Act
	service, status := validateServiceAliases(mockedClient, registration)
	// Assert
	autopilot.Equals(t, (*opslevel.Service)(nil), service)
	autopilot.Equals(t, serviceAliasesResult_NoAliasesMatched, status)
}

func Test_ValidateServiceAliases_WhenSuccessfulMatch(t *testing.T) {
	// Arrange
	mockedResponse := FixtureMockResponse{
		Status: http.StatusOK,
		Path:   "service",
	}
	mockedClient, mockedServer := AMockedClient(mockedResponse, mockedResponse, mockedResponse)
	defer mockedServer.Close()
	registration := ServiceRegistration{
		Name: "Test",
		Aliases: []string{
			"Alias1",
			"Alias2",
			"Alias3",
		},
	}
	// Act
	service, status := validateServiceAliases(mockedClient, registration)
	// Assert
	autopilot.Equals(t, "XXX", service.Id)
	autopilot.Equals(t, serviceAliasesResult_AliasMatched, status)
}

func Test_ValidateServiceAliases_WhenAliasesMatchMoreThenOneService(t *testing.T) {
	// Arrange
	mockedResponse1 := FixtureMockResponse{
		Status: http.StatusOK,
		Path:   "service",
	}
	mockedResponse2 := StringMockResponse{
		Status: http.StatusOK,
		Data:   "{}",
	}
	mockedResponse3 := FixtureMockResponse{
		Status: http.StatusOK,
		Path:   "service2",
	}
	mockedClient, mockedServer := AMockedClient(mockedResponse1, mockedResponse2, mockedResponse3)
	defer mockedServer.Close()
	registration := ServiceRegistration{
		Name: "Test",
		Aliases: []string{
			"Alias1",
			"Alias2",
			"Alias3",
		},
	}
	// Act
	service, status := validateServiceAliases(mockedClient, registration)
	// Assert
	autopilot.Equals(t, (*opslevel.Service)(nil), service)
	autopilot.Equals(t, serviceAliasesResult_MultipleServicesFound, status)
}

func Test_ValidateServiceAliases_WhenHasNon200(t *testing.T) {
	// Arrange
	mockedResponse1 := StringMockResponse{
		Status: http.StatusOK,
		Data:   "{}",
	}
	mockedResponse2 := StringMockResponse{
		Status: http.StatusRequestTimeout,
		Data:   "{}",
	}
	mockedClient, mockedServer := AMockedClient(mockedResponse1, mockedResponse2, mockedResponse1)
	defer mockedServer.Close()
	registration := ServiceRegistration{
		Name: "Test",
		Aliases: []string{
			"Alias1",
			"Alias2",
			"Alias3",
		},
	}
	// Act
	service, status := validateServiceAliases(mockedClient, registration)
	// Assert
	autopilot.Equals(t, (*opslevel.Service)(nil), service)
	autopilot.Equals(t, serviceAliasesResult_APIErrorHappened, status)
}
