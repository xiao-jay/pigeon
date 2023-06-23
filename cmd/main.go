package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"pigeon/config"
	"pigeon/pkg/framework"
	"pigeon/pkg/plugins"
	"time"
)

var MsgChan chan config.Msg

func main() {
	c := cron.New()
	Config, err := config.GetConf()
	if err != nil {
		panic(err)
	}

	if err := initcheck(*Config); err != nil {
		panic(err)
	}
	MsgChan = make(chan config.Msg, 1000)
	_, err = c.AddFunc(Config.MainCron, func() {
		log.Println("主程序开始执行任务")
		if err := SendMessage(Config.SendKeys, MsgChan); err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Println(err)
	}

	for _, crontask := range Config.CronTasks {
		log.Println(crontask)
		_, err = c.AddFunc(crontask.Cron, func() {
			MsgChan <- crontask.Message
		})
		if err != nil {
			log.Println(err)
		}
	}

	plugins.InitPlugin()
	for pluginName, argument := range Config.Plugins {
		if pluginBuilder, ok := framework.GetPluginBuilder(pluginName); ok {
			plugin := pluginBuilder(argument)
			if err := plugin.Run(MsgChan, *Config, c); err != nil {
				log.Println(err)
			}
		} else {
			log.Printf("plugin %s not found", pluginName)
		}
	}

	c.Start()
	for {
		time.Sleep(time.Second)
	}
}

func initcheck(config config.Config) error {
	if len(config.SendKeys) == 0 {
		return errors.New("no sendkey")
	}
	return nil
}

func SendMessage(sendKeys []string, msgChan chan config.Msg) error {
	if len(sendKeys) == 0 {
		return errors.New("no sendkey")
	}
	client := &http.Client{}
	msgs := make([]config.Msg, 0)
	NoneMsgFlag := false
	for !NoneMsgFlag {
		select {
		case msg := <-msgChan:
			msgs = append(msgs, msg)
		default:
			NoneMsgFlag = true
		}
	}
	if len(msgs) == 0 {
		log.Println(errors.New("no msg"))
		return nil
	}

	bytesData, _ := json.Marshal(msgs)
	for _, sendkey := range sendKeys {
		url := fmt.Sprintf("https://sctapi.ftqq.com/%s.send?title=%v&desp=%v&channel=%v", sendkey, msgs[0].Title, msgs[0].Description, msgs[0].Channel)
		req, _ := http.NewRequest("GET", url, bytes.NewReader(bytesData))
		resp, _ := client.Do(req)
		resp.Body.Close()
	}

	return nil
}
