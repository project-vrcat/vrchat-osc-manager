package config

import (
	"strings"
	"vrchat-osc-manager/internal/utils"

	"github.com/BurntSushi/toml"
)

type (
	Config struct {
		WebSocket struct {
			Hostname string `toml:"hostname"`
			Port     int    `toml:"port"`
		}
		OSC struct {
			ServerAddr string `toml:"server_addr"`
			ClientPort int    `toml:"client_port"`
		}
		RuntimePath string            `toml:"runtime_path"`
		Plugins     map[string]Plugin `toml:"plugins"`
	}
	Plugin map[string]any
)

var C Config

func LoadConfig(path string) (*Config, error) {
	_, err := toml.DecodeFile(path, &C)
	if err != nil {
		return nil, err
	}
	if C.WebSocket.Port == -1 {
		C.WebSocket.Port = utils.PickPort()
	}
	if C.RuntimePath == "" {
		C.RuntimePath = "./runtime"
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

func (p Plugin) Options() map[string]any {
	m := make(map[string]interface{})
	for k, v := range p {
		if k == "enabled" || k == "avatar_bind" {
			continue
		}
		m[k] = v
	}
	return m
}

func (p Plugin) AvatarBind(avatar string) bool {
	var avatars []string
	if _ab, ok := p["avatar_bind"]; !ok {
		return false
	} else {
		for _, a := range _ab.([]any) {
			avatars = append(avatars, a.(string))
		}
	}
	if contains(avatars, "all") {
		return true
	}
	if strings.Index(avatar, "local") == 0 && contains(avatars, "all:local") {
		return true
	}
	return contains(avatars, avatar)
}

func contains[T string](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
