package config

import "github.com/BurntSushi/toml"

type (
	Config struct {
		WebSocket struct {
			Hostname string `toml:"hostname"`
			Port     int    `toml:"port"`
		}
		Plugins map[string]Plugin `toml:"plugins"`
	}
	Plugin struct {
		Enabled bool `toml:"enabled"`
	}
)

var C Config

func LoadConfig(path string) {
	_, err := toml.DecodeFile(path, &C)
	if err != nil {
		panic(err)
	}
}
