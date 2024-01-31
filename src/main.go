package main

import "github.com/opslevel/kubectl-opslevel/cmd"

var (
	version = "dev"
	commit  = "none"
)

func main() {
	cmd.Execute(commit, version)
}
