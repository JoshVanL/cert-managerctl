package request

import (
	"crypto/x509"
	"crypto/x509/pkix"

	//"github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"

	"github.com/joshvanl/cert-managerctl/cmd/options"
	"github.com/joshvanl/cert-managerctl/pkg/util"
)

func Cert(options *options.Cert) error {

	uris, err := util.ParseURIs(options.URIs)
	if err != nil {
		return err
	}

	sk, err := util.ParsePrivateKeyFile(options.Key)
	if err != nil {
		return err
	}

	_ = x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   options.CommonName,
			Organization: options.Organization,
		},
		DNSNames:    options.DNSNames,
		IPAddresses: util.ParseIPAddresses(options.IPs),
		URIs:        uris,
		PublicKey:   sk,
	}

	return nil
}
