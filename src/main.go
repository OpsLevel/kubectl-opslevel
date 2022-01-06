package main

import (
	"fmt"

	"github.com/opslevel/kubectl-opslevel/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	cmd.Execute(fmt.Sprintf("%s-%s-%s-%s", version, commit, date, builtBy))
}
