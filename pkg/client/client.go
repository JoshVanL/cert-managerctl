package client

import (
	"io/ioutil"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	restConfig *rest.Config
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
