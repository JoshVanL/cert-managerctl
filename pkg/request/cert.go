package request

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"

	"github.com/joshvanl/cert-managerctl/cmd/options"
	"github.com/joshvanl/cert-managerctl/pkg/client"
	"github.com/joshvanl/cert-managerctl/pkg/util"
)

func Cert(client *client.Client, options *options.Cert) error {
	uris, err := util.ParseURIs(options.URIs)
	if err != nil {
		return err
	}

	// TODO: generate key if empty
	keyBundle, err := util.ParsePrivateKeyFile(options.Key)
	if err != nil {
		return err
	}

	duration, err := util.DefaultCertDuration(options.CRSpec.Duration)
	if err != nil {
		return err
	}

	csr := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   options.CommonName,
			Organization: options.Organization,
		},
		DNSNames:           options.DNSNames,
		IPAddresses:        util.ParseIPAddresses(options.IPs),
		URIs:               uris,
		PublicKey:          keyBundle.PrivateKey.Public(),
		PublicKeyAlgorithm: keyBundle.PublicKeyAlgorithm,
		SignatureAlgorithm: keyBundle.SignatureAlgorithm,
	}

	csrDER, err := util.EncodeCSR(csr, keyBundle.PrivateKey)
	if err != nil {
		return err
	}

	csrPEM := pem.EncodeToMemory(&pem.Block{
		Type: "CERTIFICATE REQUEST", Bytes: csrDER,
	})

	cr := &v1alpha1.CertificateRequest{
		ObjectMeta: util.DefaultGenerateObjectMeta(options.Object),
		Spec: v1alpha1.CertificateRequestSpec{
			CSRPEM:   csrPEM,
			IsCA:     options.CRSpec.IsCA,
			Duration: duration,
		},
	}

	err = client.CreateCertificateRequest(cr)
	if err != nil {
		return err
	}

	cr, err = client.WaitForCertificateRequestReady(
		cr.Name, cr.Namespace, time.Second*30)
	if err != nil {
		return err
	}

	fmt.Printf("Got cert:\n%s\n", cr.Status.Certificate)

	return nil
}
