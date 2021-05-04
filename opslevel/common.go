package opslevel

import (
	"fmt"
	"strings"

	"github.com/shurcooL/graphql"
)

type PageInfo struct {
	HasNextPage     graphql.Boolean `graphql:"hasNextPage"`
	HasPreviousPage graphql.Boolean `graphql:"hasPreviousPage"`
	Start           graphql.String  `graphql:"startCursor"`
	End             graphql.String  `graphql:"endCursor"`
}

type PayloadVariables map[string]interface{}

type OpsLevelErrors struct {
	Message graphql.String
	Path    []graphql.String
}

func FormatErrors(errs []OpsLevelErrors) error {
	if len(errs) == 0 {
		return nil
	}

	var errstrings []string
	errstrings = append(errstrings, "OpsLevel API Errors:")
	for _, err := range errs {
		errstrings = append(errstrings, fmt.Sprintf("\t* %s", string(err.Message)))
	}

	return fmt.Errorf(strings.Join(errstrings, "\n"))
}
