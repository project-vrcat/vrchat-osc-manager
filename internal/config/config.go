package config

import (
	"reflect"
	"strings"
	"sync"
	"vrchat-osc-manager/internal/utils"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"golang.org/x/exp/slices"
)

type (
	Config struct {
		WebSocket struct {
			Hostname string `koanf:"hostname"`
			Port     int    `koanf:"port"`
		} `koanf:"websocket"`
		OSC struct {
			ServerAddr string `koanf:"server_addr"`
			ClientPort int    `koanf:"client_port"`
		} `koanf:"osc"`
		RuntimeDir string             `koanf:"runtime_dir"`
		PluginsDir string             `koanf:"plugins_dir"`
		Plugins    map[string]*Plugin `koanf:"plugins"`
	}
	Plugin struct {
		m                  sync.Mutex
		name               string
		Enabled            bool     `koanf:"enabled"`
		AvatarBind         []string `koanf:"avatar_bind"`
		ListenAvatarChange bool
		listenParameters   []string
	}
)

var (
	C          Config
	k          = koanf.New(".")
	pluginTags []string
)

func init() {
	t := reflect.TypeOf(Plugin{})
	num := t.NumField()
	for i := 0; i < num; i++ {
		f := t.Field(i)
		tag := f.Tag.Get("koanf")
		if tag != "" && f.IsExported() {
			pluginTags = append(pluginTags, tag)
		}
	}
}

func LoadConfig(path string) (*Config, error) {
	_ = k.Load(confmap.Provider(map[string]any{
		"plugins_dir":        "./plugins",
		"runtime_dir":        "./runtime",
		"websocket.hostname": "localhost",
		"websocket.port":     utils.PickPort(),
		"osc.server_addr":    "localhost:9001",
		"osc.client_port":    9000,
	}, "."), nil)
	if err := k.Load(file.Provider(path), toml.Parser()); err != nil {
		return nil, err
	}
	if err := k.Unmarshal("", &C); err != nil {
		return nil, err
	}
	for key, plugin := range C.Plugins {
		plugin.name = key
	}
	return &C, nil
}

func (p *Plugin) SetListenParameters(params []string) {
	p.m.Lock()
	p.listenParameters = params
	p.m.Unlock()
}

func (p *Plugin) GetListenParameters() (params []string) {
	p.m.Lock()
	params = append([]string{}, p.listenParameters...)
	p.m.Unlock()
	return
}

func (p *Plugin) Options() map[string]any {
	m := k.Get("plugins." + p.name).(map[string]any)
	r := make(map[string]any)
	for k, v := range m {
		if !slices.Contains(pluginTags, k) {
			r[k] = v
		}
	}
	return r
}

func (p *Plugin) CheckAvatarBind(avatar string) (_avatar string, ok bool) {
	if p.AvatarBind == nil {
		return
	}
	if strings.Index(avatar, "local") == 0 {
		_avatar = avatar
	}
	if _avatar != "" && slices.Contains(p.AvatarBind, "all:local") {
		return _avatar, true
	}
	if slices.Contains(p.AvatarBind, "all") {
		return "all", true
	}
	ok = slices.Contains(p.AvatarBind, avatar)
	if ok {
		_avatar = avatar
	}
	return _avatar, ok
}
