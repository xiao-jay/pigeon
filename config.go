package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Plugin string

type Config struct {
	CronTasks []CronTask `yaml:"crontasks,omitempty"`
	SendKeys  []string   `yaml:"sendkeys,omitempty"`
	//Plugins   []interface{}   `yaml:"plugins"`
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
	yamlFile, err := ioutil.ReadFile("config.yaml")
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
