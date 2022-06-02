package app

import (
	"flag"
	"log"
	"vrchat-osc-manager/internal/config"
	"vrchat-osc-manager/internal/w32"
)

var (
	configFile = flag.String("config", "config.toml", "config file")
	noGUI      = flag.Bool("nogui", false, "starts the server without a gui")
	debugMode  = flag.Bool("debug", false, "enables debug mode")
)

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
	if *noGUI {
		log.Fatal(wsServer.Listen())
	} else {
		go func() { log.Fatal(wsServer.Listen()) }()
		w32.HideConsoleWindow()
		GUI()
	}
}
