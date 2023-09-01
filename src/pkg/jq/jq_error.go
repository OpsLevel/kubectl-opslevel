package jq

import (
	"fmt"
	"strings"
)

type JQErrorType int

const (
	EmptyFilter JQErrorType = iota
	BadOptions
	BadFilter
	BadJSON
	BadExcution
	ExecutableNotFound
	UnknownError
)

type JQError struct {
	Message string
	Type    JQErrorType
}

func (e JQError) Error() string {
	switch e.Type {
	case EmptyFilter:
		return "Empty JQ Filter"
	case BadOptions:
		return fmt.Sprintf("Invalid JQ Options %s", e.Message)
	case BadFilter:
		return fmt.Sprintf("Invalid JQ Filter %s", e.Message)
	case BadJSON:
		return fmt.Sprintf("Invalid Json %s", e.Message)
	case BadExcution:
		return fmt.Sprintf("Failed JQ Execution %s", strings.TrimSuffix(e.Message, "\n"))
	case ExecutableNotFound:
		return fmt.Sprintf("Executable Not Found %s", e.Message)
	case UnknownError:
		return fmt.Sprintf("Unknown JQ Error %s", e.Message)
	}
	panic(fmt.Sprintf("Unknown JQ Error %s", e.Message))
}
