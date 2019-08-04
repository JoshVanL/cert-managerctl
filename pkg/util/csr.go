package util

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/joshvanl/cert-managerctl/cmd/options"
)

// EncodeCSR calls x509.CreateCertificateRequest to sign the given CSR template.
// It returns a DER encoded signed CSR.
func EncodeCSR(template *x509.CertificateRequest, key crypto.Signer) ([]byte, error) {
	derBytes, err := x509.CreateCertificateRequest(rand.Reader, template, key)
	if err != nil {
		return nil, fmt.Errorf("error creating x509 certificate: %s", err.Error())
	}

	return derBytes, nil
}

func DefaultGenerateObjectMeta(opts options.Object) metav1.ObjectMeta {
	if len(opts.Name) == 0 {
		return metav1.ObjectMeta{
			GenerateName: "cert-managerctl-",
			Namespace:    opts.Namespace,
		}
	}

	return metav1.ObjectMeta{
		GenerateName: opts.Name,
		Namespace:    opts.Namespace,
	}
}

func DefaultCertDuration(d string) (*metav1.Duration, error) {
	if len(d) == 0 {
		return &metav1.Duration{
			Duration: v1alpha1.DefaultCertificateDuration,
		}, nil
	}

	dur, err := time.ParseDuration(d)
	if err != nil {
		return nil, err
	}

	return &metav1.Duration{
		Duration: dur,
	}, nil
}
