package app

import (
	"github.com/hypebeast/go-osc/osc"
	"log"
)

var oscClient = osc.NewClient("127.0.0.1", 9000)

func oscServer() {
	d := osc.NewStandardDispatcher()
	_ = d.AddMsgHandler("/avatar/change", func(msg *osc.Message) {
		if len(msg.Arguments) > 0 {
			log.Println("AvatarChange:", msg.Arguments[0])
		}
	})
	server := &osc.Server{
		Addr:       "127.0.0.1:9001",
		Dispatcher: d,
	}
	log.Fatal(server.ListenAndServe())
}