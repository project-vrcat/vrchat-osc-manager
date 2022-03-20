package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"vrchat-osc-manager/internal/config"
	"vrchat-osc-manager/internal/plugin"
)

const pluginsDir = "plugins"

var plugins = make(map[string]*plugin.Plugin)

func loadPlugins() {
	dir, err := ioutil.ReadDir(pluginsDir)
	if err != nil {
		panic(err)
	}
	for _, info := range dir {
		if !info.IsDir() {
			continue
		}
		enabled := false
		c, ok := config.C.Plugins[info.Name()]
		if ok {
			enabled = c.Enabled()
		}
		p, err := plugin.NewPlugin(
			filepath.Join(pluginsDir, info.Name()),
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
		if enabled {
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
