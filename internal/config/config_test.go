package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c, err := LoadConfig("../../config.toml")
	assert.Nil(t, err)

	plugin, ok := c.Plugins["example-plugin"]
	assert.Equalf(t, ok, true, "Plugin example-plugin not found")

	enabled := plugin.Enabled
	assert.Equalf(t, enabled, true, "Plugin example-plugin should be disabled")

	options := plugin.Options()
	assert.NotNil(t, options)
	assert.Equal(t, options["parameters"], []any{"VRCEmote"})

	assert.Equal(t, plugin.CheckAvatarBind("all:local"), true)
}
