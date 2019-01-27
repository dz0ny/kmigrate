package utils

import (
	"os"

	"kmigrate/logger"

	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var log = logger.New("kmigrate.utils")

func buildOutOfClusterConfig() (*rest.Config, error) {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

// GetClient returns a k8s clientset for base group
func GetClient() *rest.RESTClient {
	config, err := buildOutOfClusterConfig()
	if err != nil {
		log.Warnf("Can not get kubernetes config: %v", err)

		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("Can not get kubernetes config: %v", err)
		}
	}
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	c, err := rest.UnversionedRESTClientFor(config)

	if err != nil {
		log.Fatalf("Cannot create REST client, try setting KUBECONFIG environment variable: %v", err)
	}
	return c
}
