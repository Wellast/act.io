package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type IConf struct {
	KubeConfPath string `yaml:"kube_conf_path"`
	Steam        ISteam `yaml:"steam"`
}

type ISteam struct {
	Guardcode string `yaml:"guardcode"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

func GetConf(configPath string) (c IConf, err error) {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return IConf{}, err
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return IConf{}, err
	}

	return c, nil
}
