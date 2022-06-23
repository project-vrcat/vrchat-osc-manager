package app

import (
	"context"
	"strings"
	"vrchat-osc-manager/internal/config"
	"vrchat-osc-manager/internal/logger"

	"github.com/hypebeast/go-osc/osc"
)

type (
	OSC struct {
		ClientPort int
		ServerAddr string
		client     *osc.Client
		server     *osc.Server
	}
)

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
	ctx := context.Background()
	_ = d.AddMsgHandler("*", func(msg *osc.Message) {
		if strings.Index(msg.Address, "/avatar/parameters/") == 0 {
			parameterName := msg.Address[19:]
			for _, p := range config.C.Plugins {
				for _, param := range p.GetListenParameters() {
					if param == parameterName {
						hub.Publish(ctx, "parameters", &Message{
							method: "parameters",
							parameter: parameter{
								name:  parameterName,
								value: msg.Arguments,
							},
						})
					}
				}
			}
		}
	})
	_ = d.AddMsgHandler("/avatar/change", func(msg *osc.Message) {
		if len(msg.Arguments) > 0 {
			if avatar, ok := msg.Arguments[0].(string); ok {
				nowAvatar = avatar

				hub.Publish(ctx, "avatar_change", &Message{
					method: "avatar_change",
					avatar: nowAvatar,
				})
			}
		}
	})
	receiver.server = &osc.Server{Addr: receiver.ServerAddr, Dispatcher: d}
	logger.App("OSCServer").Println("Listening on", receiver.ServerAddr)
	return receiver.server.ListenAndServe()
}
