package cmd

import (
	"github.com/spf13/cobra"
)

var requestCertCmd = &cobra.Command{
	Use:     "certicate",
	Short:   "Request a signed certificate from cert-manager",
	Aliases: []string{"cert"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	requestCertFlags(requestCertCmd.PersistentFlags())
	requestCmd.AddCommand(requestCertCmd)
}
