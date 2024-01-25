package opslevel_k8s_controller

import (
	"encoding/json"
	"strconv"

	opslevel_jq_parser "github.com/opslevel/opslevel-jq-parser/v2024"
)

type K8SFilter struct {
	parser *opslevel_jq_parser.JQArrayParser
}

func NewK8SFilter(selector K8SSelector) *K8SFilter {
	return &K8SFilter{
		parser: opslevel_jq_parser.NewJQArrayParser(selector.Excludes),
	}
}

func (f *K8SFilter) Matches(data any) bool {
	j, err := json.Marshal(data)
	if err != nil {
		return false
	}
	// TODO: handle error
	results, _ := f.parser.Run(string(j))
	return anyIsTrue(results)
}

func anyIsTrue(results []string) bool {
	for _, result := range results {
		boolValue, err := strconv.ParseBool(result)
		if err != nil {
			return false // TODO: is this a good idea?
		}
		if boolValue {
			return true
		}
	}
	return false
}
