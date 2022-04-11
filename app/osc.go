package app

import (
	"github.com/hypebeast/go-osc/osc"
	"log"
)

type OSC struct {
	ClientPort int
	ServerAddr string
	client     *osc.Client
	server     *osc.Server
}

var nowAvatar string

func NewOSC(clientPort int, serverAddr string) *OSC {
	o := OSC{
		ClientPort: clientPort,
		ServerAddr: serverAddr,
	}
	o.client = osc.NewClient("localhost", o.ClientPort)
	return &o
}

func (receiver *OSC) Send(packet osc.Packet) error {
	return receiver.client.Send(packet)
}

func (receiver *OSC) Listen() error {
	d := osc.NewStandardDispatcher()
	_ = d.AddMsgHandler("/avatar/change", func(msg *osc.Message) {
		if len(msg.Arguments) > 0 {
			if avatar, ok := msg.Arguments[0].(string); ok {
				nowAvatar = avatar
			}
		}
	})
	receiver.server = &osc.Server{Addr: receiver.ServerAddr, Dispatcher: d}
	log.Printf("[OSCServer] Listening on %s", receiver.ServerAddr)
	return receiver.server.ListenAndServe()
}
