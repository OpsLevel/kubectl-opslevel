package opslevel

import (
	"encoding/json"
	"strconv"

	"github.com/rs/zerolog/log"
)

// JSON is a specialized map[string]string to support proper graphql serialization
type (
	JSON       map[string]any
	JSONSchema map[string]any
)

func (s JSONSchema) GetGraphQLType() string { return "JSONSchema" }

func NewJSONSchema(data string) (*JSONSchema, error) {
	result := make(JSONSchema)
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s JSONSchema) AsString() string {
	dto := map[string]any{}
	for k, v := range s {
		dto[k] = v
	}
	b, _ := json.Marshal(dto)
	return string(b)
}

func (s JSONSchema) MarshalJSON() ([]byte, error) {
	dto := map[string]any{}
	for k, v := range s {
		dto[k] = v
	}
	b, err := json.Marshal(dto)
	return []byte(strconv.Quote(string(b))), err
}

func (s JSON) GetGraphQLType() string { return "JSON" }

func NewJSON(data string) (*JSON, error) {
	result := make(JSON)
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s JSON) ToJSON() string {
	dto := map[string]any{}
	for k, v := range s {
		dto[k] = v
	}
	b, _ := json.Marshal(dto)
	return string(b)
}

func (s JSON) MarshalJSON() ([]byte, error) {
	dto := map[string]any{}
	for k, v := range s {
		dto[k] = v
	}
	b, err := json.Marshal(dto)
	return []byte(strconv.Quote(string(b))), err
}

// JsonString is a specialized input type to support serialization to JSON for input to graphql
type JsonString string

func (s JsonString) GetGraphQLType() string { return "JsonString" }

func NewJSONInput(data any) (*JsonString, error) {
	var result JsonString
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	result = JsonString(bytes)
	return &result, nil
}

func JsonStringAs[T any](data JsonString) (T, error) {
	var result T
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		log.Warn().Err(err).Msgf("unable to marshal json as %T", result)
		return result, err
	}
	return result, nil
}

func (s JsonString) AsBool() bool {
	value, _ := JsonStringAs[bool](s)
	return value
}

func (s JsonString) AsInt() int {
	value, _ := JsonStringAs[int](s)
	return value
}

func (s JsonString) AsFloat64() float64 {
	value, _ := JsonStringAs[float64](s)
	return value
}

func (s JsonString) AsString() string {
	value, _ := JsonStringAs[string](s)
	return value
}

func (s JsonString) AsArray() []any {
	value, _ := JsonStringAs[[]any](s)
	return value
}

func (s JsonString) AsMap() map[string]any {
	value, _ := JsonStringAs[map[string]any](s)
	return value
}
