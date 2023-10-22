package main

import (
	"controller/k8s"
	"flag"
	//	"fmt"

	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func main() {
	kubeconfig := *flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	k8sClient, err := k8s.GetClient(kubeconfig)
	if err != nil {
		panic(err)
	}
}
