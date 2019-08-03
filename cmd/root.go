package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/joshvanl/cert-managerctl/cmd/options"
)

var globalFlags = &options.Flags{}

var RootCmd = &cobra.Command{
	Use:   "cert-managerctl",
	Short: "A tool to interact with cert-manager via the Kubernetes API server.",
}

func Execute(args []string) {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(
		&globalFlags.Kubeconfig,
		"kubeconfig",
		"",
		"Path location to the kubeconfig file. If empty, try in-cluster authentication.",
	)
}
