package plugins

import (
	"pigeon/pkg/framework"
	"pigeon/pkg/plugins/bond"
	"pigeon/pkg/plugins/bondReminder"
	"pigeon/pkg/plugins/taopiaopiao"
	"pigeon/pkg/plugins/weather"
)

func InitPlugin() {
	framework.RegisterPluginBuilder(weather.Name, weather.New)
	framework.RegisterPluginBuilder(bondReminder.Name, bondReminder.New)
	framework.RegisterPluginBuilder(taopiaopiao.Name, taopiaopiao.New)
	framework.RegisterPluginBuilder(bond.Name, bond.New)
}
