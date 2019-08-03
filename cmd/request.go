package cmd

import (
	"github.com/jetstack/cert-manager/pkg/apis/certmanager"
	"github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/joshvanl/cert-managerctl/cmd/options"
)

var requestCmd = &cobra.Command{
	Use:     "request",
	Short:   "Request opertions on cert-manager",
	Aliases: []string{"req"},
}

func requestIssuerFlags(store *options.Issuer, fs *pflag.FlagSet) {
	fs.StringVar(
		&store.IssuerName,
		"issuer-name",
		"",
		"The target issuer name to issuer the certificate.",
	)

	fs.StringVar(
		&store.IssuerKind,
		"issuer-kind",
		"Issuer",
		"The target issuer kind to sign the certificate.",
	)

	fs.StringVar(
		&store.IssuerGroup,
		"issuer-group",
		certmanager.GroupName,
		"The target API group name the issuer belongs to",
	)
}

func requestCRSpecFlags(store *options.CRSpec, fs *pflag.FlagSet) {
	fs.StringVar(
		&store.Duration,
		"duration",
		v1alpha1.DefaultCertificateDuration.String(),
		"The requested duration the certificate will be valid for.",
	)

	fs.BoolVar(
		&store.IsCA,
		"is-ca",
		false,
		"The signed certifcate will be marked as a CA.",
	)

	fs.StringVar(
		&store.Out,
		"out",
		"/etc/cert-manager/cert.pem",
		"The output file location to store the signed certificate. If empty, output to Stdout.",
	)
}

func requestCertFlags(fs *pflag.FlagSet) {
	store := &globalFlags.Request.Cert

	fs.StringVar(
		&store.CommonName,
		"common-name",
		"",
		"Common name of the signed certificate. If empty, the first element of dns names will be used.",
	)

	fs.StringSliceVar(
		&store.Organization,
		"organisation",
		[]string{},
		"List of organisations of the signed certificate.",
	)

	fs.StringSliceVar(
		&store.DNSNames,
		"dns-names",
		[]string{},
		"List of DNS names the certificate will be valid for.",
	)

	fs.StringSliceVar(
		&store.IPs,
		"ips",
		[]string{},
		"List of IPs the certificate will be valid for.",
	)

	fs.StringSliceVar(
		&store.URIs,
		"uris",
		[]string{},
		"List of URIs the certificate will be valid for.",
	)

	fs.StringVar(
		&store.Key,
		"key",
		"/etc/cert-manager/key.pem",
		"The input key file location used to generate the CSR. If unset, an RSA "+
			"2048 private key will be generated and stored at this location",
	)

	requestCRSpecFlags(&store.CRSpec, fs)
	requestIssuerFlags(&store.Issuer, fs)
}

func requestSignFlags(fs *pflag.FlagSet) {
	store := &globalFlags.Request.Sign

	fs.StringVar(
		&store.CSR,
		"csr",
		"",
		"Path location to the CSR PEM to be signed.",
	)

	requestCRSpecFlags(&store.CRSpec, fs)
	requestIssuerFlags(&store.Issuer, fs)
}

func init() {
	RootCmd.AddCommand(requestCmd)
}
