package common

import (
	"encoding/json"

	"github.com/opslevel/kubectl-opslevel/jq"

	_ "github.com/rs/zerolog/log"
)

type JQParser struct {
	JQ jq.JQ
}

type JQResponseType int

const (
	Empty JQResponseType = iota
	String
	StringArray
	StringStringMap
	StringStringMapArray
	Unknown
)

type JQResponse struct {
	Bytes          []byte
	Type           JQResponseType
	StringObj      string
	StringArray    []string
	StringMap      map[string]string
	StringMapArray []map[string]string
}

func NewJQParser(filter string) JQParser {
	parser := JQParser{JQ: jq.New(filter)}
	return parser
}

func (parser *JQParser) doParse(data []byte) *JQResponse {
	if parser.JQ.Filter() == "" {
		return &JQResponse{Bytes: []byte("")}
	}
	var bytes []byte
	var err *jq.JQError
	bytes, err = parser.JQ.Run(data)
	if err != nil {
		//fmt.Println(err.Error())
		switch err.Type {
		case jq.BadOptions:
			return nil
		case jq.BadFilter:
			return &JQResponse{Bytes: []byte(parser.JQ.Filter())}
		case jq.BadJSON:
			return nil
		case jq.BadExcution:
			return nil
		}
	}
	return &JQResponse{Bytes: bytes}
}

func (parser *JQParser) Parse(data []byte) *JQResponse {
	resp := parser.doParse(data)
	if resp != nil {
		resp.Unmarshal()
	}
	return resp
}

func (resp *JQResponse) Unmarshal() {
	if string(resp.Bytes) == "" {
		resp.Type = Empty
		return
	}

	stringObjErr := json.Unmarshal(resp.Bytes, &resp.StringObj)
	if stringObjErr == nil {
		resp.Type = String
		return
	}

	stringArrayErr := json.Unmarshal(resp.Bytes, &resp.StringArray)
	if stringArrayErr == nil {
		resp.Type = StringArray
		return
	}

	stringMapErr := json.Unmarshal(resp.Bytes, &resp.StringMap)
	if stringMapErr == nil {
		resp.Type = StringStringMap
		return
	}

	stringMapArrayErr := json.Unmarshal(resp.Bytes, &resp.StringMapArray)
	if stringMapArrayErr == nil {
		resp.Type = StringStringMapArray
		return
	}

	resp.Type = Unknown
}
