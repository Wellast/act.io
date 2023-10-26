package main

import (
	"controller/config"
	"controller/gin"
	"controller/k8s"
	"github.com/alexflint/go-arg"
)

var args struct {
	Config string `arg:"required"`
}

func main() {
	err := arg.Parse(&args)
	if err != nil {
		panic(err)
	}

	conf, err := config.GetConf(args.Config)
	if err != nil {
		panic(err)
	}

	err = k8s.ConnectToKubernetes(conf)
	if err != nil {
		panic(err)
	}

	panic(gin.RunGin(conf))
}
