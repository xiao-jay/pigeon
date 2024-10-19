package main

import (
	"errors"
	"github.com/robfig/cron/v3"
	"log"
	"pigeon/config"
	"pigeon/pkg/channel"
	"pigeon/pkg/framework"
	"pigeon/pkg/plugins"
	"pigeon/pkg/router"
)

var MsgChan chan config.Msg

func main() {
	// default print file name and line
	log.SetFlags(log.Lshortfile)
	c := cron.New()
	Config, err := config.GetConf("config/config.yaml")
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

	go router.InitRouter(Config)
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
	f := channels.FangTang{}
	if err := f.SendMessage(msgs, sendKeys); err != nil {
		return err
	}
	return nil
}
