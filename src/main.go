package main

import (
	"fmt"

	"github.com/opslevel/kubectl-opslevel/cmd"
)

var (
	version = "dev"
	commit  = "none"
)

func truncate_commit(commit string, length int) (short_commit string) {
	if len(commit) > length {
		return commit[:length]
	}
	return commit
}

func main() {
	cmd.Execute(fmt.Sprintf("%s-%s", version, truncate_commit(commit, 12)))
}
