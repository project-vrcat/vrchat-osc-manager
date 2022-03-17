package config

import (
	"github.com/BurntSushi/toml"
)

type (
	Config struct {
		WebSocket struct {
			Hostname string `toml:"hostname"`
			Port     int    `toml:"port"`
		}
		Plugins map[string]Plugin `toml:"plugins"`
	}
	Plugin map[string]interface{}
)

var C Config

func LoadConfig(path string) (*Config, error) {
	_, err := toml.DecodeFile(path, &C)
	if err != nil {
		return nil, err
	}
	return &C, nil
}

func (p Plugin) Enabled() bool {
	e, ok := p["enabled"]
	if ok {
		if e, ok := e.(bool); ok {
			return e
		}
	}
	return false
}

func (p Plugin) Options() map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range p {
		if k == "enabled" {
			continue
		}
		m[k] = v
	}
	return m
}
