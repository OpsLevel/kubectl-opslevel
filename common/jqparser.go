package common

import (
	"fmt"
	"encoding/json"

	"github.com/opslevel/kubectl-opslevel/jq"

	_ "github.com/rs/zerolog/log"
)

type JQParser struct {
	JQ jq.JQ
}

func NewJQParser(filter string) JQParser {
	parser := JQParser{JQ: jq.New(filter)}
	return parser
}

func (parser *JQParser) Parse(data []byte) string {
	if (parser.JQ.Filter() == "") {
		return ""
	}
	var bytes []byte
	var err *jq.JQError
	bytes, err = parser.JQ.Run(data)
	if err != nil {
		//fmt.Println(err.Error())
		switch err.Type {
		case jq.BadOptions:
			return ""
		case jq.BadFilter:
			return parser.JQ.Filter()
		case jq.BadJSON:
			return ""
		case jq.BadExcution:
			return ""
		}
	}
	var output string
	jsonErr := json.Unmarshal(bytes, &output)
	if (jsonErr != nil) {
		// TODO: `jq all` returns a "bool" - guess we should handle this?
		// TODO: Unable to marshal to string so we have a jq return of some other type - [] or map[string]interface
		fmt.Printf("Failed To Parse JQ Return `%s`: %s\n", parser.JQ.Commandline(), jsonErr.Error())
		return ""
	}
	return output
}