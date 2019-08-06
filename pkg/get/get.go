package get

import (
	"fmt"
	"time"

	cmapi "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"

	"github.com/joshvanl/cert-managerctl/cmd/options"
	"github.com/joshvanl/cert-managerctl/pkg/client"
	"github.com/joshvanl/cert-managerctl/pkg/util"
)

type Get struct {
	client *client.Client
	opts   *options.Get
}

func New(client *client.Client, opts *options.Get) *Get {
	return &Get{
		client: client,
		opts:   opts,
	}
}

func (g *Get) Cert() error {
	opts := g.opts.Cert

	cr, err := g.getOrWait(opts.Wait, &opts.Object)
	if err != nil {
		return err
	}

	if !util.CertificateRequestReady(cr) {
		return fmt.Errorf("certificate request %s/%s not ready: %s: %s",
			cr.Name, cr.Namespace, cr.Status.Conditions)
	}

	if out := opts.Out; len(out) > 0 {
		return util.WriteCertificateFile(out, cr.Status.Certificate)
	}

	fmt.Printf("%s", cr.Status.Certificate)

	return nil
}

func (g *Get) getOrWait(wait bool, opts *options.Object) (*cmapi.CertificateRequest, error) {
	if !wait {
		return g.client.CertificateRequest(opts.Name, opts.Namespace)
	}

	return g.client.WaitForCertificateRequestReady(opts.Name, opts.Namespace, time.Second*30)
}
