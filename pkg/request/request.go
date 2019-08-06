package request

import (
	"fmt"
	"time"

	cmapi "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	log "github.com/sirupsen/logrus"

	"github.com/joshvanl/cert-managerctl/cmd/options"
	"github.com/joshvanl/cert-managerctl/pkg/client"
	"github.com/joshvanl/cert-managerctl/pkg/util"
)

type Request struct {
	client *client.Client
	opts   *options.Request
}

func New(client *client.Client, opts *options.Request) *Request {
	return &Request{
		client: client,
		opts:   opts,
	}
}

func (r *Request) csr(csrPEM []byte, opts *options.CROptions) error {
	duration, err := util.DefaultCertDuration(opts.CRSpec.Duration)
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

	log.Info("creating CertificateRequest")

	cr, err = r.client.CreateCertificateRequest(cr)
	if err != nil {
		return err
	}

	log.Infof("waiting for CertificateRequest %s/%s to become ready",
		cr.Namespace, cr.Name)
	cr, err = r.client.WaitForCertificateRequestReady(
		cr.Name, cr.Namespace, time.Second*30)
	if err != nil {
		return fmt.Errorf("failed waiting for resource %s/%s to become ready: %s",
			cr.Namespace, cr.Name, err)
	}

	log.Info("signed certificate successfully issued")

	if out := opts.CRSpec.Out; len(out) > 0 {
		log.Infof("writing signed certificate request to %s", out)

		return util.WriteFile(out, cr.Status.Certificate, 0600)
	}

	fmt.Printf("%s", cr.Status.Certificate)

	return nil
}
