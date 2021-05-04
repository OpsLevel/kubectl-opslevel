package jq

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
)

type JQ struct {
	options []string
	timeout time.Duration
	writer  io.Writer
}

type JQOpt struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type JQErrorType int

const (
	EmptyFilter JQErrorType = iota
	BadOptions
	BadFilter
	BadJSON
	BadExcution
	Unknown
)

type JQError struct {
	Message string
	Type    JQErrorType
}

func (e *JQError) Error() string {
	switch e.Type {
	case EmptyFilter:
		return "Empty JQ Filter"
	case BadOptions:
		return fmt.Sprintf("Bad JQ Options: `%s`", e.Message)
	case BadFilter:
		return fmt.Sprintf("Bad JQ Filter: `%s`", e.Message)
	case BadJSON:
		return fmt.Sprintf("Bad Json: `%s`", e.Message)
	case BadExcution:
		return fmt.Sprintf("Bad JQ Execution: `%s`", e.Message)
	case Unknown:
		return fmt.Sprintf("Unknown JQ Error: `%s`", e.Message)
	}
	panic(fmt.Sprintf("Unknown JQ Error Occured!"))
}

func (jq *JQ) Filter() string {
	return jq.options[len(jq.options)-1]
}

func (jq *JQ) Options() []string {
	return jq.options[:len(jq.options)-1]
}

func (jq *JQ) Commandline() string {
	return fmt.Sprintf("jq %s", strings.Join(jq.options, " "))
}

func (jq *JQ) Run(json []byte) ([]byte, *JQError) {
	var stderr bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), jq.timeout)
	//fmt.Printf("Exec: `jq %s`\n", strings.Join(jq.options, " "))
	cmd := exec.CommandContext(ctx, "jq", jq.options...)
	cmd.Stdin = bytes.NewBuffer(json)
	cmd.Stderr = &stderr
	cmd.Env = make([]string, 0)
	defer cancel()
	out, err := cmd.Output()
	if err != nil {
		//fmt.Println("Got Error on JQ Execution")
		//fmt.Println(err.Error())
		//fmt.Println(string(stderr.Bytes()))
		if err.Error() == "exit status 2" {
			return nil, &JQError{Message: jq.Commandline(), Type: BadOptions}
		}
		if err.Error() == "exit status 3" {
			return nil, &JQError{Message: jq.Filter(), Type: BadFilter}
		}
		if err.Error() == "exit status 4" {
			return nil, &JQError{Message: string(json), Type: BadJSON}
		}
		if err.Error() == "exit status 5" {
			return nil, &JQError{Message: string(stderr.Bytes()), Type: BadExcution}
		}
		return nil, &JQError{Message: string(stderr.Bytes()), Type: BadExcution}
	}
	return out, nil
}

func (jq *JQ) Validate(json []byte) *JQError {
	filter := jq.Filter()
	if filter == "" {
		return &JQError{Message: filter, Type: EmptyFilter}
	}
	_, err := jq.Run(json)
	return err
}

func New(filter string) JQ {
	return NewWithOptions(filter, 8*time.Second, nil)
}

func NewWithOptions(filter string, timeout time.Duration, options []JQOpt) JQ {
	opts := []string{}
	for _, opt := range options {
		if opt.Enabled {
			opts = append(opts, fmt.Sprintf("--%s", opt.Name))
		}
	}
	opts = append(opts, fmt.Sprintf("%s", filter))
	jq := &JQ{
		options: opts,
		timeout: timeout,
		writer:  ioutil.Discard,
	}
	return *jq
}
