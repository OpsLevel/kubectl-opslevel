package opslevel

import (
	"fmt"
	"strings"

	"github.com/shurcooL/graphql"
)

type PayloadVariables map[string]interface{}

type Lifecycle struct {
	Alias graphql.String
	Description graphql.String
	Id graphql.ID
	Index graphql.Int
	Name graphql.String
}

type Tier struct {
	Alias graphql.String
	Description graphql.String
	Id graphql.ID
	Index graphql.Int
	Name graphql.String
}

type OpsLevelErrors struct {
	Message graphql.String
	Path []graphql.String
}

func FormatErrors(errs []OpsLevelErrors) error {
	if (len(errs) == 0) {
		return nil
	}

	var errstrings []string 
	errstrings = append(errstrings, "OpsLevel API Errors:")
	for _, err := range errs {
		errstrings = append(errstrings, fmt.Sprintf("\t* %s", string(err.Message)))
	}
	
	return fmt.Errorf(strings.Join(errstrings, "\n"))
}