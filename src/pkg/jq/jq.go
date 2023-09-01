package jq

import (
	"bytes"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"os/exec"
	"strings"
	"time"
)

var _validated = false

type JQ struct {
	binary    string
	options   []string
	filter    string
	timeout   time.Duration
	writer    io.Writer
	validated bool
}

func NewJQ() *JQ {
	return &JQ{
		binary:  "jq",
		filter:  ".",
		options: []string{},
		timeout: 10 * time.Second,
		writer:  io.Discard,
	}
}

func (jq *JQ) validate() bool {
	if jq.validated {
		return true
	}
	_, err := exec.LookPath(jq.binary)
	if err != nil {
		return false
	}
	jq.validated = true
	return true
}

func (jq *JQ) WithBinary(binary string) *JQ {
	jq.binary = binary
	return jq
}

func (jq *JQ) WithOption(option string) *JQ {
	jq.options = append(jq.options, option)
	return jq
}

func (jq *JQ) Raw() *JQ {
	return jq.WithOption("-r")
}

func (jq *JQ) WithFilter(filter string) *JQ {
	jq.filter = filter
	return jq
}

func (jq *JQ) WithTimeout(timeout time.Duration) *JQ {
	jq.timeout = timeout
	return jq
}

func (jq *JQ) commandline() string {
	return fmt.Sprintf("%s %s %s", jq.binary, strings.Join(jq.options, " "), jq.filter)
}

func (jq *JQ) run(ctx context.Context, data []byte, stderr *bytes.Buffer) ([]byte, error) {
	log.Debug().Msgf("Exec: '%s'", jq.commandline())
	cmd := exec.CommandContext(ctx, jq.binary, append(jq.options, jq.filter)...)
	cmd.Stdin = bytes.NewBuffer(data)
	cmd.Stderr = stderr
	cmd.Env = make([]string, 0)
	return cmd.Output()
}

func (jq *JQ) Exec(data []byte) (*JQResponse, error) {
	if jq.validate() == false {
		//log.Fatal().Msgf("Please install '%s' to use this tool - https://stedolan.github.io/jq/download/", jq.binary)
		return nil, JQError{Message: jq.binary, Type: ExecutableNotFound}
	}
	if jq.filter == "" {
		return nil, JQError{Message: jq.filter, Type: EmptyFilter}
	}
	var stderr bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), jq.timeout)
	defer cancel()
	out, err := jq.run(ctx, data, &stderr)
	if err != nil {
		log.Debug().Err(err).Str("stderr", stderr.String()).Msg("error occurred while executing JQ")
		if err.Error() == "exit status 2" {
			return nil, JQError{Message: jq.commandline(), Type: BadOptions}
		}
		if err.Error() == "exit status 3" {
			return nil, JQError{Message: jq.filter, Type: BadFilter}
		}
		if err.Error() == "exit status 4" {
			return nil, JQError{Message: string(data), Type: BadJSON}
		}
		if err.Error() == "exit status 5" {
			return nil, JQError{Message: stderr.String(), Type: BadExcution}
		}
		return nil, JQError{Message: stderr.String(), Type: BadExcution}
	}
	return NewResponse(out), nil
}

func (jq *JQ) Run(data string) (*JQResponse, error) {
	return jq.Exec([]byte(data))
}
