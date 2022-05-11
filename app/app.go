package app

import (
	"flag"
	"log"
	"vrchat-osc-manager/internal/config"
)

var configFile = flag.String("config", "config.toml", "config file")

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Start() {
	flag.Parse()

	if _, err := config.LoadConfig(*configFile); err != nil {
		log.Fatal(err)
	}

	osc := NewOSC(config.C.OSC.ClientPort, config.C.OSC.ServerAddr)
	wsServer := NewWSServer(config.C.WebSocket.Hostname, config.C.WebSocket.Port, osc)

	go func() { log.Fatal(osc.Listen(wsServer)) }()
	go loadPlugins()
	go func() { log.Fatal(wsServer.Listen()) }()
	GUI()
}
