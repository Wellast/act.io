package k8s

import (
	"controller/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var k8sClient *kubernetes.Clientset
var appConfig config.IConf

func ConnectToKubernetes(conf config.IConf) error {
	appConfig = conf

	restConf, err := clientcmd.BuildConfigFromFlags("", appConfig.KubeConfPath)
	if err != nil {
		return err
	}

	k8sClient, err = kubernetes.NewForConfig(restConf)
	if err != nil {
		return err
	}

	return nil
}
