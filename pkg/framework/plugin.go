package framework

import (
	"github.com/robfig/cron/v3"
	"pigeon/config"
)

func RunPlugin(c *cron.Cron) error {

	return nil
}

// PluginBuilder plugin management
type PluginBuilder = func(config.Arguments) Plugin

// Plugin management
var pluginBuilders = map[string]PluginBuilder{}

// RegisterPluginBuilder register the plugin
func RegisterPluginBuilder(name string, pc PluginBuilder) {
	pluginBuilders[name] = pc
}

// GetPluginBuilder get the pluginbuilder by name
func GetPluginBuilder(name string) (PluginBuilder, bool) {
	pb, found := pluginBuilders[name]
	return pb, found
}
