package common

import (
	"encoding/json"
	"fmt"

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
	Bool
	BoolArray
	Unknown
)

type JQResponse struct {
	Bytes          []byte
	Type           JQResponseType
	StringObj      string
	StringArray    []string
	StringMap      map[string]string
	StringMapArray []map[string]string
	BoolObj        bool
	BoolArray      []bool
}

type JQResponseMulti struct {
	Bytes   []byte
	Objects []JQResponse
}

func NewJQParser(filter string) JQParser {
	parser := JQParser{JQ: jq.New(filter)}
	return parser
}

func NewJQParserMulti(filter string) JQParser {
	parser := JQParser{JQ: jq.New(fmt.Sprintf("map(%s)", filter))}
	return parser
}

func (parser *JQParser) doParse(data []byte) []byte {
	var bytes []byte
	var err *jq.JQError
	bytes, err = parser.JQ.Run(data)
	if err != nil {
		//fmt.Println(err.Error())
		switch err.Type {
		case jq.BadOptions:
			return nil
		case jq.BadFilter:
			return []byte(fmt.Sprintf("\"%s\"", parser.JQ.Filter()))
		case jq.BadJSON:
			return nil
		case jq.BadExcution:
			return nil
		}
	}
	return bytes
}

func (parser *JQParser) Parse(data []byte) *JQResponse {
	var resp *JQResponse
	if parser.JQ.Filter() == "" {
		resp = &JQResponse{Bytes: []byte("")}
	} else {
		resp = &JQResponse{Bytes: parser.doParse(data)}
	}
	resp.Unmarshal()
	return resp
}

func (parser *JQParser) ParseMulti(data []byte) *JQResponseMulti {
	var resp *JQResponseMulti
	if parser.JQ.Filter() == "map()" {
		resp = &JQResponseMulti{Bytes: []byte("[]")}
	} else {
		resp = &JQResponseMulti{Bytes: parser.doParse(data)}
	}
	resp.Unmarshal()
	return resp
}

func (resp *JQResponse) Unmarshal() {
	//fmt.Printf("Unmarshaling '%s'\n", string(resp.Bytes))
	if string(resp.Bytes) == "" {
		resp.Type = Empty
		return
	}

	stringObjErr := json.Unmarshal(resp.Bytes, &resp.StringObj)
	if stringObjErr == nil {
		if resp.StringObj == "" {
			resp.Type = Empty
			return
		}
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

	boolObjErr := json.Unmarshal(resp.Bytes, &resp.BoolObj)
	if boolObjErr == nil {
		resp.Type = Bool
		return
	}

	boolArrayErr := json.Unmarshal(resp.Bytes, &resp.BoolArray)
	if boolArrayErr == nil {
		resp.Type = BoolArray
		return
	}

	resp.Type = Unknown
}

func (resp *JQResponseMulti) Unmarshal() {
	//fmt.Printf("Unmarshaling '%s'\n", string(resp.Bytes))
	if string(resp.Bytes) == "[]" {
		resp.Objects = nil
		return
	}

	var multi_stringObj []string
	var multi_stringArray [][]string
	var multi_stringMap []map[string]string
	var multi_stringMapArray [][]map[string]string

	stringObjErr := json.Unmarshal(resp.Bytes, &multi_stringObj)
	if stringObjErr == nil {
		for _, item := range multi_stringObj {
			resp.Objects = append(resp.Objects, JQResponse{Type: String, StringObj: item})
		}
		return
	}

	stringArrayErr := json.Unmarshal(resp.Bytes, &multi_stringArray)
	if stringArrayErr == nil {
		for _, item := range multi_stringArray {
			resp.Objects = append(resp.Objects, JQResponse{Type: StringArray, StringArray: item})
		}
		return
	}

	stringMapErr := json.Unmarshal(resp.Bytes, &multi_stringMap)
	if stringMapErr == nil {
		for _, item := range multi_stringMap {
			resp.Objects = append(resp.Objects, JQResponse{Type: StringStringMap, StringMap: item})
		}
		return
	}

	stringMapArrayErr := json.Unmarshal(resp.Bytes, &multi_stringMapArray)
	if stringMapArrayErr == nil {
		for _, item := range multi_stringMapArray {
			resp.Objects = append(resp.Objects, JQResponse{Type: StringStringMapArray, StringMapArray: item})
		}
		return
	}

	resp.Objects = nil
}
