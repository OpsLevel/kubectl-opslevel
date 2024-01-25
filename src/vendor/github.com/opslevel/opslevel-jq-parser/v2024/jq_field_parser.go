package opslevel_jq_parser

import (
	"fmt"

	"github.com/flant/libjq-go"
	"github.com/flant/libjq-go/pkg/jq"
)

type JQFieldParser struct {
	program *jq.JqProgram
}

func NewJQFieldParser(expression string) *JQFieldParser {
	if expression == "" {
		expression = "empty"
	}
	prg, err := libjq_go.Jq().Program(expression).Precompile()
	if err != nil {
		panic(fmt.Sprintf("unable to compile jq expression:  %s", expression))
	}
	return &JQFieldParser{
		program: prg,
	}
}

func (p *JQFieldParser) Run(data string) (string, error) {
	return p.program.RunRaw(data)
}
