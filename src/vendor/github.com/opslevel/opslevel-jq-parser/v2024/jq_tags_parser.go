package opslevel_jq_parser

import (
	"encoding/json"
	"strings"

	"github.com/opslevel/opslevel-go/v2024"
)

type JQTagsParser struct {
	creates []*JQFieldParser
	assigns []*JQFieldParser
}

func NewJQTagsParser(cfg TagRegistrationConfig) *JQTagsParser {
	creates := make([]*JQFieldParser, len(cfg.Create))
	for i, expression := range cfg.Create {
		creates[i] = NewJQFieldParser(expression)
	}
	assigns := make([]*JQFieldParser, len(cfg.Assign))
	for i, expression := range cfg.Assign {
		assigns[i] = NewJQFieldParser(expression)
	}
	return &JQTagsParser{
		creates: creates,
		assigns: assigns,
	}
}

func (p *JQTagsParser) parse(programs []*JQFieldParser, data string) []opslevel.TagInput {
	output := make([]opslevel.TagInput, 0, len(programs))
	for _, program := range programs {
		response, err := program.Run(data)
		if err != nil {
			// TODO: log error
			continue
		}
		if response == "" {
			continue
		}

		if strings.HasPrefix(response, "[") && strings.HasSuffix(response, "]") {
			var tags []map[string]string
			if err := json.Unmarshal([]byte(response), &tags); err != nil {
				// TODO: log error
				continue
			}
			for _, item := range tags {
				for key, value := range item {
					if key == "" || value == "" {
						// TODO: log warning
						continue
					}
					output = append(output, opslevel.TagInput{Key: key, Value: value})
				}
			}
		}
		if strings.HasPrefix(response, "{") && strings.HasSuffix(response, "}") {
			var tags map[string]string
			if err := json.Unmarshal([]byte(response), &tags); err != nil {
				// TODO: log error
				continue
			}
			for key, value := range tags {
				if key == "" || value == "" {
					// TODO: log warning
					continue
				}
				output = append(output, opslevel.TagInput{Key: key, Value: value})
			}
		}
	}
	return output
}

func (p *JQTagsParser) Run(data string) ([]opslevel.TagInput, []opslevel.TagInput, error) {
	return p.parse(p.creates, data), p.parse(p.assigns, data), nil
}
