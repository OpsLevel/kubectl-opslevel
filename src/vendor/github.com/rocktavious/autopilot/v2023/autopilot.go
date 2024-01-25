package autopilot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

/*
var (
	Client *MyClient
)

func TestMain(m *testing.M) {
	flag.Parse()
	teardown := autopilot.Setup()
	defer teardown()
	Client := NewClient(autopilot.Server.URL)
	os.Exit(m.Run())
}

func TestExample(t *testing.T) {
	// Arrange
	autopilot.Mux.HandleFunc("/orgs/octokit/repos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, autopilot.Fixture("repos/octokit.json"))
	})
	// Act
	body, err := Client.DoStuff()
	a := Client.GetThing("Four")
	// Assert
	autopilot.Assert(t, a >= 4, "value must be greater than or equal to four")
	autopilot.Ok(t, err)
	autopilot.Equals(t, []byte("OK"), body)
}
*/

var (
	Mux       *http.ServeMux
	Server    *httptest.Server
	Templater *FixtureTemplater
)

// Setup an HttpTestServer and ServerMux (for path routing) and return the teardown function
func Setup() func() {
	Mux = http.NewServeMux()
	Server = httptest.NewServer(Mux)
	Templater = NewFixtureTemplater()
	return func() {
		Server.Close()
	}
}

type RequestValidation func(*http.Request)

type ResponseWriter func(http.ResponseWriter)

func SkipRequestValidation() RequestValidation {
	return func(r *http.Request) {}
}

func EmptyResponse() ResponseWriter {
	return func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "{}")
	}
}

func JsonStringResponse(jsonString string) ResponseWriter {
	if !json.Valid([]byte(jsonString)) {
		panic(fmt.Errorf("invalid json passed to JsonStringResponse(): %s", jsonString))
	}
	return func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, jsonString)
	}
}

func FixtureResponse(fixture string) ResponseWriter {
	response := TemplatedFixture(fixture)
	return func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, response)
	}
}

func RegisterEndpoint(endpoint string, responseWriter ResponseWriter, requestValidation RequestValidation) string {
	Mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		requestValidation(r)
		responseWriter(w)
	})
	return Server.URL + endpoint
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// Returns a Mock'd HttpClient
func NewHttpClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

/*
client := autopilot.NewHttpClient(func(req *http.Request) *http.Response {
	// Test request parameters
	autopilot.Equals(t, req.URL.String(), "http://example.com/some/path")
	return &http.Response{
		StatusCode: 200,
		Body: io.NopCloser(bytes.NewBufferString(`OK`)),
		// Must be set to non-nil value or it panics
		Header: make(http.Header),
	}
})
api := API{client, "http://example.com"}
*/

// Load testdata/fixtures/<path> and return as string
func Fixture(path string) string {
	b, err := os.ReadFile("testdata/fixtures/" + path)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func TemplatedFixture(fixture string) string {
	response, err := Templater.Use(Fixture(fixture))
	if err != nil {
		panic(err)
	}
	return response
}

// RunTableTests runs a table of tests calling `fn` for each test and ensuring all tests run in parallel using goroutines
func RunTableTests[T any](t *testing.T, cases map[string]T, fn func(*testing.T, T)) {
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			fn(t, test)
		})
	}
}

// Assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		tb.Helper()
		tb.Fatalf(msg, v...)
	}
}

// Ok fails the test if an err is not nil.
func Ok(tb testing.TB, err error) {
	if err != nil {
		tb.Helper()
		tb.Fatalf("\n\tunexpected error: %s", err.Error())
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		tb.Helper()
		tb.Fatalf("\n\texp: %#v\n\tgot: %#v", exp, act)
	}
}

func Register[T any](name string, value T) T {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	_, err = Templater.coreTemplate.Parse(fmt.Sprintf(`{{ define "%s" }}
%s
{{ end }}`, name, data))
	if err != nil {
		panic(err)
	}
	return value
}
