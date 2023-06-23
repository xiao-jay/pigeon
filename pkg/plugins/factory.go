package plugins

import (
	"pigeon/pkg/framework"
	"pigeon/pkg/plugins/weather"
)

func InitPlugin() {
	framework.RegisterPluginBuilder(weather.Name, weather.New)
}
