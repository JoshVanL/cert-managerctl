package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "v0.1.0-alpha1"

var versionCmd = &cobra.Command{
	Args:  cobra.NoArgs,
	Use:   "version",
	Short: "prints the cert-managerctl CLI version",
	Long:  "prints the cert-managerctl CLI version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(Version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
