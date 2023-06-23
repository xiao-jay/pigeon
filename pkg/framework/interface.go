package framework

import (
	"github.com/robfig/cron/v3"
	"pigeon/config"
)

type Plugin interface {
	// Name The unique name of Plugin.
	Name() string
	// SendMessage send msg to channel
	SendMessage(msg interface{}, Msg chan config.Msg) error
	Run(Msg chan config.Msg, config config.Config, c *cron.Cron) error
}
