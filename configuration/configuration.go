package configuration

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type DataBaseConfig struct {
	ID       string `json:"ID"`
	Password string `json:"Password"`
	Name     string `json:"Name"`
}

type SMSConfig struct {
	SecretId   string `json:"SecretId"`
	SecretKey  string `json:"SecretKey"`
	SdkAppId   string `json:"SdkAppId"`
	Signature  string `json:"Signature"`
	TemplateId string `json:"TemplateId"`
}

type config struct {
	Port         string         `json:"Port"`
	ViewRootPath string         `json:"ViewRootPath"`
	EEGDataPath  string         `json:"EEGDataPath"`
	DB           DataBaseConfig `json:"DB"`
	SMS          SMSConfig      `json:"SMS"`
}

func (c *config) init() {
	work_dir, _ := os.Getwd()
	bytes, err := os.ReadFile(work_dir + "/configuration/configuration.json")
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(bytes, c)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%+v\n", c)
}

var configInstance *config
var configOnce sync.Once

func GetConfigInstance() *config {
	configOnce.Do(func() {
		configInstance = &config{}
		configInstance.init()
	})
	return configInstance
}
