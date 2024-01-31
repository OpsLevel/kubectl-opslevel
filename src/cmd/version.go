package cmd

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var commit, version string

type Build struct {
	Version string `json:"version,omitempty"`
	Commit  string `json:"git_commit,omitempty"`
	GoInfo  GoInfo `json:"go,omitempty"`
}

type GoInfo struct {
	Version  string `json:"version,omitempty"`
	Compiler string `json:"compiler,omitempty"`
	OS       string `json:"os,omitempty"`
	Arch     string `json:"arch,omitempty"`
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print version information`,
	RunE:  runVersion,
}

func runVersion(cmd *cobra.Command, args []string) error {
	goInfo := GoInfo{
		Version:  runtime.Version(),
		Compiler: runtime.Compiler,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
	}
	build := Build{Commit: commit, GoInfo: goInfo, Version: version}
	versionInfo, err := json.MarshalIndent(build, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(versionInfo))
	return nil
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
