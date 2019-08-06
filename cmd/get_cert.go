package cmd

import (
	"github.com/joshvanl/cert-managerctl/pkg/client"
	"github.com/joshvanl/cert-managerctl/pkg/get"
	"github.com/spf13/cobra"
)

var getCertCmd = &cobra.Command{
	Use:     "certicate",
	Short:   "Get the certificate stored in a CertificateRequest.",
	Aliases: []string{"cert"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := client.New(flags.Kubeconfig)
		if err != nil {
			return err
		}

		get := get.New(client, &flags.Get)
		mustDie(get.Cert())

		return nil
	},
}

func init() {
	getCertFlags(getCertCmd.PersistentFlags())
	getCmd.AddCommand(getCertCmd)
}
