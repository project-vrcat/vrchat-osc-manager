package app

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"vrchat-osc-manager/internal/config"
	"vrchat-osc-manager/internal/plugin"
)

func init() {
	config.LoadConfig("./config.toml")
}

func Start() {
	go oscServer()
	go loadPlugins()
	wsServer()
}

func loadPlugins() {
	dir, err := ioutil.ReadDir("./plugins/")
	if err != nil {
		panic(err)
	}
	for _, info := range dir {
		if !info.IsDir() {
			continue
		}
		p, err := plugin.NewPlugin(filepath.Join("./plugins/", info.Name()))
		if err != nil {
			log.Println(err)
			continue
		}
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
