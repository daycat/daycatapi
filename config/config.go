package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var ApiKey string
var ApiEmail string
var Zoneid string

type configFile struct {
	ApiKey   string `yaml:"ApiKey"`
	ApiEmail string `yaml:"ApiEmail"`
	Zoneid   string `yaml:"Zoneid"`
}

func GetConfig(path string) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var config configFile
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	ApiKey = config.ApiKey
	ApiEmail = config.ApiEmail
	Zoneid = config.Zoneid

}
