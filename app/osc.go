package app

import (
	"strings"
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
	parameterInfo struct {
		Name  string
		Value []any
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

func (receiver *OSC) Listen(ws *WSServer) error {
	d := osc.NewStandardDispatcher()
	_ = d.AddMsgHandler("*", func(msg *osc.Message) {
		if strings.Index(msg.Address, "/avatar/parameters/") == 0 {
			parameterName := msg.Address[19:]
			pluginsParameters.Range(func(k, v any) bool {
				pluginName := k.(string)
				parameters := v.([]string)
				for _, parameter := range parameters {
					if parameter == parameterName {
						if ch, ok := ws.parameterChan.Load(pluginName); ok {
							ch.(chan parameterInfo) <- parameterInfo{
								Name:  parameterName,
								Value: msg.Arguments,
							}
						}
					}
				}
				return true
			})
		}
	})
	_ = d.AddMsgHandler("/avatar/change", func(msg *osc.Message) {
		if len(msg.Arguments) > 0 {
			if avatar, ok := msg.Arguments[0].(string); ok {
				nowAvatar = avatar
				pluginsAvatarChange.Range(func(k, v any) bool {
					pluginName := k.(string)
					if ch, ok := ws.avatarChangeChan.Load(pluginName); ok {
						ch.(chan bool) <- true
					}
					return true
				})
			}
		}
	})
	receiver.server = &osc.Server{Addr: receiver.ServerAddr, Dispatcher: d}
	logger.App("OSCServer").Println("Listening on", receiver.ServerAddr)
	return receiver.server.ListenAndServe()
}
