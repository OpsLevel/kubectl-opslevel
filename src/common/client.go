package common

import (
	"fmt"
	"strings"

	"github.com/opslevel/opslevel-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewClient() *opslevel.Client {
	client := opslevel.NewClient(viper.GetString("apitoken"), opslevel.SetURL(viper.GetString("apiurl")))

	clientErr := client.Validate()
	if clientErr != nil {
		if strings.Contains(clientErr.Error(), "Please provide a valid OpsLevel API token") {
			cobra.CheckErr(fmt.Errorf("%s via 'export OL_APITOKEN=XXX' or '--api-token=XXX' or '--api-token-path=/path/to/token/file'", clientErr.Error()))
		} else {
			cobra.CheckErr(clientErr)
		}
	}
	cobra.CheckErr(clientErr)

	return client
}
