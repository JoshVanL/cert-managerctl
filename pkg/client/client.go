package client

import (
	"fmt"
	"io/ioutil"
	"time"

	apiutil "github.com/jetstack/cert-manager/pkg/api/util"
	"github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
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
	var restConfig *rest.Config
	var err error

	if len(kubeconfig) == 0 {
		restConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}

	} else {
		kubeconfigBytes, err := ioutil.ReadFile(kubeconfig)
		if err != nil {
			return nil, err
		}

		restConfig, err = clientcmd.RESTConfigFromKubeConfig(kubeconfigBytes)
		if err != nil {
			return nil, err
		}
	}

	return &Client{
		restConfig: restConfig,
	}, nil
}

func (c *Client) CreateCertificateRequest(cr *v1alpha1.CertificateRequest) error {
	_, err := c.cmClient.CertmanagerV1alpha1().CertificateRequests(cr.Namespace).Create(cr)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) WaitForCertificateRequestReady(name, ns string, timeout time.Duration) (*v1alpha1.CertificateRequest, error) {
	var cr *v1alpha1.CertificateRequest
	err := wait.PollImmediate(time.Second, timeout,
		func() (bool, error) {
			var err error
			//log.Logf("Waiting for CertificateRequest %s to be ready", name)
			cr, err = c.cmClient.CertmanagerV1alpha1().CertificateRequests(ns).Get(name, metav1.GetOptions{})
			if err != nil {
				return false, fmt.Errorf("error getting CertificateRequest %s: %v", name, err)
			}
			isReady := apiutil.CertificateRequestHasCondition(cr, v1alpha1.CertificateRequestCondition{
				Type:   v1alpha1.CertificateRequestConditionReady,
				Status: v1alpha1.ConditionTrue,
			})
			if !isReady {
				//log.Logf("Expected CertificateReques to have Ready condition 'true' but it has: %v", cr.Status.Conditions)
				return false, nil
			}
			return true, nil
		},
	)

	if err != nil {
		return nil, err
	}

	return cr, nil
}
