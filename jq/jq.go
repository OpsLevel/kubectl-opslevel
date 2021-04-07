package jq

import (
	"context"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"time"
	"io/ioutil"
)

type JQ struct {
	options []string
	timeout time.Duration
	writer io.Writer
}

type JQOpt struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

func (jq *JQ) Run(json []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), jq.timeout)
	cmd := exec.CommandContext(ctx, "jq", jq.options...)
	cmd.Stdin = bytes.NewBuffer(json)
	cmd.Env = make([]string, 0)
	defer cancel()
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func Create(filter string) (JQ, error) {
	return CreateWithOptions(filter, 8*time.Second, nil)
}

func CreateWithOptions(filter string, timeout time.Duration, options []JQOpt) (JQ, error) {
	opts := []string{}
	for _, opt := range options {
		if opt.Enabled {
			opts = append(opts, fmt.Sprintf("--%s", opt.Name))
		}
	}
	opts = append(opts, filter)
	jq := &JQ{
		options: opts,
		timeout: timeout,
		writer: ioutil.Discard,
	}
	return *jq, nil
}