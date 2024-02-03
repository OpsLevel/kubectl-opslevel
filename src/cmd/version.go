package cmd

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/opslevel/opslevel-go/v2024"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	commit, version  string
	shortVersionFlag bool
)

type Build struct {
	Version         string          `json:"version,omitempty"`
	Commit          string          `json:"git_commit,omitempty"`
	GoInfo          GoInfo          `json:"go,omitempty"`
	OpslevelVersion OpslevelVersion `json:"opslevel,omitempty"`
}

type OpslevelVersion struct {
	Commit  string `json:"app_commit,omitempty"`
	Version string `json:"app_version,omitempty"`
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
	if shortVersionFlag {
		fmt.Printf("%s-%s\n", version, commit)
		return nil
	}
	goInfo := GoInfo{
		Version:  runtime.Version(),
		Compiler: runtime.Compiler,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
	}
	build := Build{Commit: commit, GoInfo: goInfo, Version: version}
	build.OpslevelVersion = getOpslevelVersion()
	versionInfo, err := json.MarshalIndent(build, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(versionInfo))
	return nil
}

func getOpslevelVersion() OpslevelVersion {
	opslevelVersion := OpslevelVersion{}
	clientRest := opslevel.NewRestClient(opslevel.SetURL(viper.GetString("api-url")))
	_, err := clientRest.R().SetResult(&opslevelVersion).Get("api/ping")
	cobra.CheckErr(err)

	if len(opslevelVersion.Commit) >= 12 {
		opslevelVersion.Commit = opslevelVersion.Commit[:12]
	}

	return opslevelVersion
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.PersistentFlags().BoolVar(&shortVersionFlag, "short", false, "Print only version number")
}
