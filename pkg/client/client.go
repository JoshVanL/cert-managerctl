package client

import (
	"fmt"
	"io/ioutil"
	"time"

	// This package is required to be imported to register all client
	// plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	cmapi "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	cmclient "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/joshvanl/cert-managerctl/pkg/util"
)

type Client struct {
	restConfig *rest.Config
	cmClient   cmclient.Interface
}

func New(kubeconfig string) (*Client, error) {
	restConfig, err := restConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	cmClient, err := cmclient.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return &Client{
		restConfig: restConfig,
		cmClient:   cmClient,
	}, nil
}

func (c *Client) CreateCertificateRequest(
	cr *cmapi.CertificateRequest) (*cmapi.CertificateRequest, error) {
	return c.cmClient.CertmanagerV1alpha1().CertificateRequests(cr.Namespace).Create(cr)
}

func (c *Client) CertificateRequest(name, ns string) (*cmapi.CertificateRequest, error) {
	return c.cmClient.CertmanagerV1alpha1().CertificateRequests(ns).Get(name, metav1.GetOptions{})
}

func (c *Client) WaitForCertificateRequestReady(name, ns string, timeout time.Duration) (*cmapi.CertificateRequest, error) {
	var cr *cmapi.CertificateRequest
	err := wait.PollImmediate(time.Second, timeout,
		func() (bool, error) {

			log.Debugf("polling CertificateRequest %s/%s for ready status", name, ns)

			var err error
			cr, err = c.cmClient.CertmanagerV1alpha1().CertificateRequests(ns).Get(name, metav1.GetOptions{})
			if err != nil {
				return false, fmt.Errorf("error getting CertificateRequest %s: %v", name, err)
			}

			if reason, failed := util.CertificateRequestFailed(cr); failed {
				return false, fmt.Errorf("certificate request marked as failed: %s", reason)
			}

			if !util.CertificateRequestReady(cr) {
				return false, nil
			}

			return true, nil
		},
	)

	if err != nil {
		return cr, err
	}

	return cr, nil
}

func restConfig(kubeconfig string) (*rest.Config, error) {
	if len(kubeconfig) == 0 {
		restConfig, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}

		return restConfig, nil
	}

	kubeconfigBytes, err := ioutil.ReadFile(kubeconfig)
	if err != nil {
		return nil, err
	}

	restConfig, err := clientcmd.RESTConfigFromKubeConfig(kubeconfigBytes)
	if err != nil {
		return nil, err
	}

	return restConfig, nil
}
