package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var k8sClient *kubernetes.Clientset

func ConnectToKubernetes(kubeConfPath string) error {

	restConf, err := clientcmd.BuildConfigFromFlags("", kubeConfPath)
	if err != nil {
		return err
	}

	k8sClient, err = kubernetes.NewForConfig(restConf)
	if err != nil {
		return err
	}

	return nil
}
