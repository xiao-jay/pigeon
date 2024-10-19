package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Arguments map[string]interface{}
type Config struct {
	CronTasks      []CronTask           `yaml:"crontasks"`
	SendKeys       []string             `yaml:"sendkeys"`
	Plugins        map[string]Arguments `yaml:"plugins"`
	MainCron       string               `yaml:"maincron"`
	FeishuWebHooks []string             `yaml:"feishu_webhooks"`
}

type CronTask struct {
	Name    string `json:"name"`
	Cron    string `json:"cron"`
	Message Msg    `json:"message"`
}

type Msg struct {
	Title       string `json:"title"`
	Description string `json:"desp"`
	Channel     int    `json:"channel"`
}

func GetConf(confg_file_path string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(confg_file_path)
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
