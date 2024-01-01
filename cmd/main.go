package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"io"
	"log"
	"net/http"
	"pigeon/config"
	"pigeon/pkg/framework"
	"pigeon/pkg/plugins"
	"strconv"
)

var MsgChan chan config.Msg

func main() {
	// default print file name and line
	log.SetFlags(log.Lshortfile)
	c := cron.New()
	Config, err := config.GetConf()
	if err != nil {
		panic(err)
	}

	if err := initcheck(*Config); err != nil {
		panic(err)
	}
	MsgChan = make(chan config.Msg, 1000)
	log.Println("主程序开始执行任务")
	_, err = c.AddFunc(Config.MainCron, func() {
		if err := SendMessage(Config.SendKeys, MsgChan); err != nil {
			log.Println("Error:", err)
		}
	})
	if err != nil {
		log.Println(err)
	}

	for _, crontask := range Config.CronTasks {
		log.Println(crontask)
		task := config.CronTask{
			Name:    crontask.Name,
			Cron:    crontask.Cron,
			Message: crontask.Message,
		}
		_, err = c.AddFunc(task.Cron, func() {
			MsgChan <- task.Message
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
	log.Printf("have %v cron job ", len(c.Entries()))
	c.Run()
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
	log.Println("send msgs", msgs)
	// send to all sendkeys people
	if err := sendMessageThroughFangTang(msgs, sendKeys); err != nil {
		return err
	}
	return nil
}

func sendMessageThroughFangTang(msgs []config.Msg, sendKeys []string) error {
	for _, sendkey := range sendKeys {
		client := &http.Client{}
		url := fmt.Sprintf("https://sctapi.ftqq.com/%s.send", sendkey)
		for _, msg := range msgs {
			jsondata, err := json.Marshal(msg)
			if err != nil {
				return err
			}
			req, err := http.NewRequest("POST", url, bytes.NewReader(jsondata))
			if err != nil {
				return err
			}
			req.Header.Add("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			if resp.StatusCode != 200 {
				return fmt.Errorf(strconv.Itoa(resp.StatusCode))
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Println(err)
				}
			}(resp.Body)
		}
	}

	return nil
}
