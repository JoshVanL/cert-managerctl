package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joshvanl/cert-managerctl/pkg/client"
	"github.com/joshvanl/cert-managerctl/pkg/request"
)

var requestCertCmd = &cobra.Command{
	Use:     "certicate",
	Short:   "Request a signed certificate from cert-manager.",
	Aliases: []string{"cert"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := client.New(flags.Kubeconfig)
		if err != nil {
			return err
		}

		request := request.New(client, &flags.Request)
		mustDie(request.Cert())

		return nil
	},
}

func init() {
	requestCertFlags(requestCertCmd.PersistentFlags())
	requestCmd.AddCommand(requestCertCmd)
}
