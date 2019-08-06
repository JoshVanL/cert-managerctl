package client

import (
	"fmt"
	"io/ioutil"
	"time"

	// This package is required to be imported to register all client
	// plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	apiutil "github.com/jetstack/cert-manager/pkg/api/util"
	cmapi "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	cmclient "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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

func (c *Client) CreateCertificateRequest(
	cr *cmapi.CertificateRequest) (*cmapi.CertificateRequest, error) {
	return c.cmClient.CertmanagerV1alpha1().CertificateRequests(cr.Namespace).Create(cr)
}

func (c *Client) WaitForCertificateRequestReady(name, ns string, timeout time.Duration) (*cmapi.CertificateRequest, error) {
	var cr *cmapi.CertificateRequest
	err := wait.PollImmediate(time.Second, timeout,
		func() (bool, error) {
			var err error
			cr, err = c.cmClient.CertmanagerV1alpha1().CertificateRequests(ns).Get(name, metav1.GetOptions{})
			if err != nil {
				return false, fmt.Errorf("error getting CertificateRequest %s: %v", name, err)
			}

			isReady := apiutil.CertificateRequestHasCondition(cr, cmapi.CertificateRequestCondition{
				Type:   cmapi.CertificateRequestConditionReady,
				Status: cmapi.ConditionTrue,
			})
			if !isReady {
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
