package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c, err := LoadConfig("../../config.toml")
	assert.Nil(t, err)

	plugin, ok := c.Plugins["example-plugin"]
	assert.Equalf(t, ok, true, "Plugin pulsoid-plugin not found")

	enabled := plugin.Enabled()
	assert.Equalf(t, enabled, true, "Plugin pulsoid-plugin should be disabled")

	options := plugin.Options()
	assert.NotNil(t, options)
	assert.Equal(t, options["parameters"], "635e7cd8-0def-4d91-a108-c198f122f663")

	assert.Equal(t, plugin.AvatarBind("all:local"), true)
}
