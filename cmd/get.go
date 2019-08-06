package cmd

import (
	"github.com/joshvanl/cert-managerctl/cmd/options"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get cert-manager resources.",
}

func getCertFlags(fs *pflag.FlagSet) {
	store := &flags.Get.Cert

	fs.StringVarP(
		&store.Out,
		"out",
		"o",
		"",
		"The output file location to store the signed certificate. If empty, output "+
			"to Stdout.",
	)

	fs.BoolVarP(
		&store.Wait,
		"wait",
		"w",
		false,
		"Wait for the target CertificateRequest to become ready",
	)

	getCertObjectFlags(&store.Object, fs)
}

func getCertObjectFlags(store *options.Object, fs *pflag.FlagSet) {
	fs.StringVar(
		&store.Name,
		"name",
		"",
		"The name of the CertificateRequest storing the certificate.",
	)

	fs.StringVarP(
		&store.Namespace,
		"namespace",
		"n",
		"",
		"The namespace of the CertificateRequest storing the certificate.",
	)
}

func init() {
	rootCmd.AddCommand(getCmd)
}
