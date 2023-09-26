package common

import (
	"encoding/json"
	"github.com/opslevel/opslevel-go/v2023"
	"github.com/rs/zerolog/log"
	"strings"
)

type RepositoryDTO struct {
	Name      string
	Directory string
	Repo      string
}

func (r *RepositoryDTO) Convert() opslevel.ServiceRepositoryCreateInput {
	return opslevel.ServiceRepositoryCreateInput{
		Repository:    *opslevel.NewIdentifier(r.Repo),
		BaseDirectory: r.Directory,
		DisplayName:   r.Name,
	}
}

type JQRepositoryParser struct {
	programs []*JQFieldParser
}

func NewJQRepositoryParser(expressions []string) *JQRepositoryParser {
	programs := make([]*JQFieldParser, len(expressions))
	for i, expression := range expressions {
		programs[i] = NewJQFieldParser(expression)
	}
	return &JQRepositoryParser{
		programs: programs,
	}
}

func (p *JQRepositoryParser) Run(data string) ([]opslevel.ServiceRepositoryCreateInput, error) {
	output := make([]opslevel.ServiceRepositoryCreateInput, 0, len(p.programs))
	for _, program := range p.programs {
		response, err := program.Run(data)
		//log.Warn().Msgf("expression: %s\nresponse: %s", program.program.Program, response)
		if err != nil {
			log.Warn().Msgf("unable to parse alias from expression: %s", program.program.Program)
			return nil, err
		}
		if response == "" {
			continue
		}
		if strings.HasPrefix(response, "[") && strings.HasSuffix(response, "]") {
			if response == "[]" {
				continue
			}
			var repos []RepositoryDTO
			if err := json.Unmarshal([]byte(response), &repos); err == nil {
				for _, repo := range repos {
					output = append(output, repo.Convert())
				}
			} else {
				// Try as []string
				var repoNames []string
				if err := json.Unmarshal([]byte(response), &repoNames); err == nil {
					for _, repoName := range repoNames {
						output = append(output, opslevel.ServiceRepositoryCreateInput{Repository: *opslevel.NewIdentifier(repoName)})
					}
				}
			}
		} else if strings.HasPrefix(response, "{") && strings.HasSuffix(response, "}") {
			var repo RepositoryDTO
			if err := json.Unmarshal([]byte(response), &repo); err != nil {
				// TODO: log error
				log.Warn().Err(err).Msgf("unable to marshal repo expression: %s\n%s", program.program.Program, response)
				continue
			}
			output = append(output, repo.Convert())
		} else {
			output = append(output, opslevel.ServiceRepositoryCreateInput{Repository: *opslevel.NewIdentifier(response)})
		}

	}
	return output, nil
}
