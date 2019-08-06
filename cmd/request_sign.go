package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joshvanl/cert-managerctl/pkg/client"
	"github.com/joshvanl/cert-managerctl/pkg/request"
)

var requestSignCmd = &cobra.Command{
	Use:   "sign",
	Short: "Request a signed certificate from cert-manager using a raw x509 encoded certificate siging request.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := client.New(flags.Kubeconfig)
		if err != nil {
			return err
		}

		request := request.New(client, &flags.Request)
		mustDie(request.Sign())

		return nil
	},
}

func init() {
	requestSignFlags(requestSignCmd.PersistentFlags())
	requestCmd.AddCommand(requestSignCmd)
}
