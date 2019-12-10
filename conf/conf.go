package conf

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"log"
)

const configFilePath string = "config-dev.yaml"

var configFile []byte

type AppConfig struct {
	App App `yaml:"app"`
}

type App struct {
	Database Database `yaml:"database"`
	Redis Redis `yaml:"redis"`
}

type Database struct {
	Type string `yaml:"type"`
	User string `yaml:"user"`
	Password string `yaml:"password"`
	DbName string `yaml:"dbName"`
	Address string `yaml:"address"`
	MaxIdle int `yaml:"maxIdle"`
	MaxOpen int `yaml:"maxOpen"`
}

type Redis struct {
	Address string `yaml:"address"`
	Network string `yaml:"network"`
	Password string `yaml:"password"`
	MaxIdle int `yaml:"maxIdle"`
	MaxActive int `yaml:"maxActive"`
	IdleTimeout int `yaml:"idleTimeout"`
}

func GetAppConfig() (appConfig AppConfig, err error)  {
	err = yaml.Unmarshal(configFile, &appConfig)
	return appConfig, err
}

func init()  {
	var err error
	configFile, err = ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("yamlFile.Get err %v", err)
	}
}