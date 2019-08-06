package request

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	cmapi "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"

	"github.com/joshvanl/cert-managerctl/cmd/options"
	"github.com/joshvanl/cert-managerctl/pkg/client"
	"github.com/joshvanl/cert-managerctl/pkg/util"
)

func Cert(client *client.Client, opts *options.Cert) error {
	uris, err := util.ParseURIs(opts.URIs)
	if err != nil {
		return err
	}

	duration, err := util.DefaultCertDuration(opts.CRSpec.Duration)
	if err != nil {
		return err
	}

	keyBundle, err := privateKey(opts.Key)
	if err != nil {
		return err
	}

	commonName, err := commonName(opts)
	if err != nil {
		return err
	}

	csr := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: opts.Organization,
		},
		DNSNames:           opts.DNSNames,
		IPAddresses:        util.ParseIPAddresses(opts.IPs),
		URIs:               uris,
		PublicKey:          keyBundle.PrivateKey.Public(),
		PublicKeyAlgorithm: keyBundle.PublicKeyAlgorithm,
		SignatureAlgorithm: keyBundle.SignatureAlgorithm,
	}

	csrPEM, err := util.EncodeCSR(csr, keyBundle.PrivateKey)
	if err != nil {
		return err
	}

	cr := &cmapi.CertificateRequest{
		ObjectMeta: util.DefaultGenerateObjectMeta(opts.Object),
		Spec: cmapi.CertificateRequestSpec{
			CSRPEM:   csrPEM,
			IsCA:     opts.CRSpec.IsCA,
			Duration: duration,
			IssuerRef: cmapi.ObjectReference{
				Name:  opts.Issuer.Name,
				Kind:  opts.Issuer.Kind,
				Group: opts.Issuer.Group,
			},
		},
	}

	cr, err = client.CreateCertificateRequest(cr)
	if err != nil {
		return err
	}

	cr, err = client.WaitForCertificateRequestReady(
		cr.Name, cr.Namespace, time.Second*30)
	if err != nil {
		return err
	}

	if out := opts.CRSpec.Out; len(out) > 0 {
		util.WriteFile(out, cr.Status.Certificate, 0600)
	} else {
		fmt.Printf("%s", cr.Status.Certificate)
	}

	return nil
}

func privateKey(keyPath string) (*util.KeyBundle, error) {
	exists, err := util.FileExists(keyPath)
	if err != nil {
		return nil, err
	}

	var keyBundle *util.KeyBundle
	if !exists {
		sk, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}

		keyPEM := pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(sk),
			},
		)

		err = util.WriteFile(keyPath, keyPEM, 0600)
		if err != nil {
			return nil, err
		}

		return &util.KeyBundle{
			PrivateKey:         sk,
			SignatureAlgorithm: x509.SHA256WithRSA,
			PublicKeyAlgorithm: x509.RSA,
		}, nil
	}

	keyBundle, err = util.ParsePrivateKeyFile(keyPath)
	if err != nil {
		return nil, err
	}

	return keyBundle, nil
}

func commonName(opts *options.Cert) (string, error) {
	if len(opts.CommonName) > 0 {
		return opts.CommonName, nil
	}

	if len(opts.DNSNames) == 0 {
		return "", errors.New("no common name or DNS names given")
	}

	return opts.DNSNames[0], nil
}
