package opslevel_jq_parser

import (
	"encoding/json"
	"strings"

	"github.com/opslevel/opslevel-go/v2024"
	"github.com/rs/zerolog/log"
)

type JQToolsParser struct {
	programs []*JQFieldParser
}

func NewJQToolsParser(expressions []string) *JQToolsParser {
	programs := make([]*JQFieldParser, len(expressions))
	for i, expression := range expressions {
		programs[i] = NewJQFieldParser(expression)
	}
	return &JQToolsParser{
		programs: programs,
	}
}

func (p *JQToolsParser) Run(data string) ([]opslevel.ToolCreateInput, error) {
	output := make([]opslevel.ToolCreateInput, 0, len(p.programs))
	for _, program := range p.programs {
		response, err := program.Run(data)
		if err != nil {
			log.Warn().Msgf("unable to parse alias from expression: %s", program.program.Program)
			return nil, err
		}
		if response == "" {
			continue
		}
		// TODO: response can be []map[string]string also
		if strings.HasPrefix(response, "[") && strings.HasSuffix(response, "]") {
			var tools []opslevel.ToolCreateInput
			if err := json.Unmarshal([]byte(response), &tools); err != nil {
				log.Err(err).Msgf("unable to parse expression: %s", program.program.Program)
				continue
			}
			output = append(output, tools...)
		} else {
			var tool opslevel.ToolCreateInput
			if err := json.Unmarshal([]byte(response), &tool); err != nil {
				log.Err(err).Msgf("unable to parse expression: %s", program.program.Program)
				continue
			}
			output = append(output, tool)
		}
	}
	return output, nil
}
