package main

import (
	"controller/gin"
	"controller/k8s"
	"github.com/alexflint/go-arg"
)

var args struct {
	Kubeconfig string `arg:"required"`
}

func main() {
	arg.MustParse(&args)

	err := k8s.ConnectToKubernetes(args.Kubeconfig)
	if err != nil {
		panic(err)
	}

	panic(gin.RunGin())
}
