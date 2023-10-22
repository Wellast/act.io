package main

import (
	"controller/gin"
	"controller/k8s"
	"flag"
	//	"fmt"

	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func main() {
	kubeconfig := *flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	err := k8s.ConnectToKubernetes(kubeconfig)
	if err != nil {
		panic(err)
	}

	panic(gin.RunGin())
}
