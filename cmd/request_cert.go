package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joshvanl/cert-managerctl/pkg/client"
	"github.com/joshvanl/cert-managerctl/pkg/request"
)

var requestCertCmd = &cobra.Command{
	Use:     "certicate",
	Short:   "Request a signed certificate from cert-manager",
	Aliases: []string{"cert"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := client.New(globalFlags.Kubeconfig)
		if err != nil {
			return err
		}

		return request.Cert(client, &globalFlags.Request.Cert)
	},
}

func init() {
	requestCertFlags(requestCertCmd.PersistentFlags())
	requestCmd.AddCommand(requestCertCmd)
}
