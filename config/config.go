package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Arguments map[string]interface{}
type Config struct {
	CronTasks []CronTask           `yaml:"crontasks"`
	SendKeys  []string             `yaml:"sendkeys"`
	Plugins   map[string]Arguments `yaml:"plugins"`
	MainCron  string               `yaml:"maincron"`
}

type CronTask struct {
	Name    string `json:"name"`
	Cron    string `json:"cron"`
	Message Msg    `json:"message"`
}

type Msg struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Channel     int    `json:"channel"`
}

func GetConf() (*Config, error) {
	yamlFile, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	log.Println("config yaml", string(yamlFile))
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	log.Printf("config:\n%+v", config)
	return &config, nil
}
