package jq_test

import (
	"github.com/opslevel/kubectl-opslevel/pkg/jq"
	"github.com/rocktavious/autopilot"
	"testing"
	"time"
)

func TestJQ(t *testing.T) {
	// Arrange
	data := `{"metadata":{"name":"test"}}`
	jq1 := jq.NewJQ().WithTimeout(1 * time.Second).WithFilter(".metadata.name")
	// Act
	result1, err1 := jq1.Run(data)
	// Assert
	autopilot.Ok(t, err1)
	autopilot.Equals(t, "test", result1.Data)
}

func TestJQRaw(t *testing.T) {
	// Arrange
	data := `{"metadata":{"name":"test"}}`
	jq1 := jq.NewJQ().Raw().WithFilter(".metadata.name")
	// Act
	result1, err1 := jq1.Run(data)
	// Assert
	autopilot.Ok(t, err1)
	autopilot.Equals(t, []byte("test\n"), result1.Data)
}

func TestJQErrors(t *testing.T) {
	// Arrange
	data := `{"metadata":{"name":"test"}}`
	bad := `{[}]`
	jq1 := jq.NewJQ().WithTimeout(1).WithFilter(".metadata.name")
	jq2 := jq.NewJQ().WithFilter("")
	jq3 := jq.NewJQ().WithOption("--p")
	jq4 := jq.NewJQ().WithFilter("..p")
	jq5 := jq.NewJQ()
	jq6 := jq.NewJQ().WithBinary("gojq")
	// Act
	_, err1 := jq1.Run(data)
	_, err2 := jq2.Run(data)
	_, err3 := jq3.Run(data)
	_, err4 := jq4.Run(data)
	_, err5 := jq5.Run(bad)
	_, err6 := jq6.Run(data)
	// Assert
	autopilot.Assert(t, err1 != nil, "Timeout Didn't Happen")
	autopilot.Assert(t, (err1.(jq.JQError)).Type == jq.BadExcution, "Error Type Was Wrong")
	autopilot.Assert(t, err2 != nil, "Didn't Handle Empty Filter")
	autopilot.Assert(t, (err2.(jq.JQError)).Type == jq.EmptyFilter, "Error Type Was Wrong")
	autopilot.Assert(t, err3 != nil, "Didn't Handle Bad Options")
	autopilot.Assert(t, (err3.(jq.JQError)).Type == jq.BadOptions, "Error Type Was Wrong")
	autopilot.Assert(t, err4 != nil, "Didn't Handle Bad Filter")
	autopilot.Assert(t, (err4.(jq.JQError)).Type == jq.BadFilter, "Error Type Was Wrong")
	autopilot.Assert(t, err5 != nil, "Didn't Handle Bad JSON")
	autopilot.Assert(t, (err5.(jq.JQError)).Type == jq.BadJSON, "Error Type Was Wrong")
	autopilot.Assert(t, err6 != nil, "Didn't Handle Invalid Binary")
	autopilot.Assert(t, (err6.(jq.JQError)).Type == jq.ExecutableNotFound, "Error Type Was Wrong")

}

func TestJQMulti(t *testing.T) {
	// Arrange
	jq1 := jq.NewJQ().WithFilter("map((.metadata.name) // null)")
	// Act
	result1, err1 := jq1.Run(`[{"metadata":{"name":"test1"}},{"metadata":{"name":"test2"}}]`)
	// Assert
	autopilot.Ok(t, err1)
	autopilot.Equals(t, []any{"test1", "test2"}, result1.Data)
}

func TestJQReponses_Empty(t *testing.T) {
	// Arrange
	data := `{"metadata":{"name":""}}`
	jq1 := jq.NewJQ().WithFilter(".metadata.name")
	jq2 := jq.NewJQ().WithFilter(".metadata.test")
	// Act
	result1, err1 := jq1.Run(data)
	result2, err2 := jq2.Run(data)
	// Assert
	autopilot.Ok(t, err1)
	autopilot.Equals(t, jq.Empty, result1.Type)
	autopilot.Equals(t, "", result1.Data)
	autopilot.Ok(t, err2)
	autopilot.Equals(t, jq.Empty, result2.Type)
	autopilot.Equals(t, "", result2.Data)
}

func TestJQReponses_Map(t *testing.T) {
	// Arrange
	jq1 := jq.NewJQ().WithFilter("{\"key\": .metadata.name}")
	// Act
	result1, err1 := jq1.Run(`{"metadata":{"name":"test2"}}`)
	// Assert
	autopilot.Ok(t, err1)
	autopilot.Equals(t, jq.Map, result1.Type)
	autopilot.Equals(t, map[string]any{"key": "test2"}, result1.Data)
}

func TestJQReponses_ArrayMap(t *testing.T) {
	// Arrange
	jq1 := jq.NewJQ().WithFilter("map({\"key\": .metadata.name} // null)")
	// Act
	result1, err1 := jq1.Run(`[{"metadata":{"name":"test1"}},{"metadata":{"name":"test2"}}]`)
	// Assert
	autopilot.Ok(t, err1)
	autopilot.Equals(t, jq.ArrayMap, result1.Type)
	autopilot.Equals(t, []map[string]any{{"key": "test1"}, {"key": "test2"}}, result1.Data)
}

func TestJQReponses_Array(t *testing.T) {
	// Arrange
	jq1 := jq.NewJQ().WithFilter("map(.metadata.name // null)")
	// Act
	result1, err1 := jq1.Run(`[{"metadata":{"name":"test1"}},{"metadata":{"name":"test2"}}]`)
	// Assert
	autopilot.Ok(t, err1)
	autopilot.Equals(t, jq.Array, result1.Type)
	autopilot.Equals(t, []any{"test1", "test2"}, result1.Data)
}
