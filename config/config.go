package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var ApiKey string
var ApiEmail string
var Zoneid string
var RootDomain string
var OutsideURL string

type configFile struct {
	ApiKey     string `yaml:"ApiKey"`
	ApiEmail   string `yaml:"ApiEmail"`
	Zoneid     string `yaml:"Zoneid"`
	RootDomain string `yaml:"RootDomain"`
	OutsideURL string `yaml:"OutsideURL"`
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
	RootDomain = config.RootDomain
	OutsideURL = config.OutsideURL

}
