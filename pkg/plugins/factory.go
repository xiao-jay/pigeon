package plugins

import (
	"pigeon/pkg/framework"
	"pigeon/pkg/plugins/bondReminder"
	"pigeon/pkg/plugins/weather"
)

func InitPlugin() {
	framework.RegisterPluginBuilder(weather.Name, weather.New)
	framework.RegisterPluginBuilder(bondReminder.Name, bondReminder.New)
}
