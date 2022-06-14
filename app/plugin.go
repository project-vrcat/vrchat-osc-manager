package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"
	"vrchat-osc-manager/internal/config"
	"vrchat-osc-manager/internal/plugin"
)

var plugins = make(map[string]*plugin.Plugin)
var pluginsParameters = sync.Map{}
var pluginsAvatarChange = sync.Map{}

func loadPlugins() {
	dir, err := ioutil.ReadDir(config.C.PluginsDir)
	if err != nil {
		panic(err)
	}
	for _, info := range dir {
		if !info.IsDir() {
			continue
		}
		c, ok := config.C.Plugins[info.Name()]
		if !ok {
			continue
		}
		p, err := plugin.NewPlugin(
			filepath.Join(config.C.PluginsDir, info.Name()),
			fmt.Sprintf("%s:%d", config.C.WebSocket.Hostname, config.C.WebSocket.Port),
		)
		if err != nil {
			log.Println(err)
			continue
		}
		if _, ok := plugins[p.Name]; ok {
			log.Printf("Plugin already loaded: %s: %s\n", p.Name, p.Dir)
			continue
		}
		plugins[p.Name] = p
		pluginsParameters.Store(p.Name, []string{})
		if c.Enabled {
			if err = p.Init(); err != nil {
				log.Println(err)
				continue
			}
			if err = p.Start(); err != nil {
				log.Println(err)
				continue
			}
		}
	}
}
