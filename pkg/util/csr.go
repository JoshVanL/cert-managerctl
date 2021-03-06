package util

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	cmapi "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/joshvanl/cert-managerctl/cmd/options"
)

// EncodeCSR calls x509.CreateCertificateRequest to sign the given CSR.
// It returns a PEM encoded signed CSR.
func EncodeCSR(csr *x509.CertificateRequest, key crypto.Signer) ([]byte, error) {
	derBytes, err := x509.CreateCertificateRequest(rand.Reader, csr, key)
	if err != nil {
		return nil, fmt.Errorf("error creating x509 certificate: %s", err.Error())
	}

	csrPEM := pem.EncodeToMemory(&pem.Block{
		Type: "CERTIFICATE REQUEST", Bytes: derBytes,
	})

	return csrPEM, nil
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
			Duration: cmapi.DefaultCertificateDuration,
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

func CertificateRequestFailed(cr *cmapi.CertificateRequest) (string, bool) {
	for _, con := range cr.Status.Conditions {
		if con.Reason == "Failed" {
			return con.Message, true
		}
	}

	return "", false
}

func CertificateRequestReady(cr *cmapi.CertificateRequest) bool {
	readyType := cmapi.CertificateRequestConditionReady
	readyStatus := cmapi.ConditionTrue

	existingConditions := cr.Status.Conditions
	for _, cond := range existingConditions {
		if readyType == cond.Type && readyStatus == cond.Status {
			return true
		}
	}

	return false
}
