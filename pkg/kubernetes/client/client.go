package client

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
)

// GetInclusterConfig returns the config on k8s cluster
func GetInclusterConfig() (kubernetes.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("getting client config for Kubernetes client: %w", err)
	}
	return kubernetes.NewForConfig(config)
}

// GetK8sRestclient returns a REST client config for API calls against the Kubernetes API for active context.
func GetK8sRestclient() (kubernetes.Interface, error) {
	rawConfig, err := getCurrentConfig()
	if err != nil {
		return nil, fmt.Errorf("getting client config for Kubernetes client: %w", err)
	}
	clientConfig := clientcmd.NewNonInteractiveClientConfig(rawConfig, "", &clientcmd.ConfigOverrides{}, clientcmd.NewDefaultClientConfigLoadingRules())
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("error creating REST client config in-cluster: %w", err)
	}
	return kubernetes.NewForConfig(restConfig)
}

func getCurrentConfig() (clientcmdapi.Config, error) {
	bytes, err := getKubeconfigBytes()
	if err != nil {
		return clientcmdapi.Config{}, fmt.Errorf("could not find kubeconfig: %w", err)

	}
	k8sConfig, err := clientcmd.NewClientConfigFromBytes(bytes)
	if err != nil {
		return clientcmdapi.Config{}, err
	}
	cfg, err := k8sConfig.RawConfig()
	return cfg, err
}

// getKubeconfigPath currently returns a single path $(HOMEDIR)/.kube/config.
// TODO (tejal29): Support adding multiple kube config files via env variable or okteto config
func getKubeconfigBytes() ([]byte, error) {
	fp := filepath.Join(homedir.HomeDir(), ".kube", "config")
	bytes, err := os.ReadFile(fp)
	if err != nil && !os.IsNotExist(err) {
		return []byte{}, err
	}
	return bytes, err
}
