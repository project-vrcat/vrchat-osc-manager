package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"
	"vrchat-osc-manager/internal/config"
	"vrchat-osc-manager/internal/logger"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/hypebeast/go-osc/osc"
)

type (
	WSServer struct {
		hostname         string
		port             int
		osc              *OSC
		parameterChan    sync.Map
		avatarChangeChan sync.Map
	}

	wsMessage struct {
		Method         string         `json:"method"`
		Plugin         string         `json:"plugin"`
		Addr           string         `json:"addr,omitempty"`
		Value          any            `json:"value,omitempty"`
		Options        map[string]any `json:"options,omitempty"`
		Parameters     []string       `json:"parameters,omitempty"`
		ParameterName  string         `json:"parameterName,omitempty"`
		ParameterValue any            `json:"parameterValue,omitempty"`
	}
)

func NewWSServer(hostname string, port int, osc *OSC) *WSServer {
	return &WSServer{
		hostname: hostname,
		port:     port,
		osc:      osc,
	}
}

func (s *WSServer) handle(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Println(err)
		return
	}
	pluginName := r.URL.Query().Get("plugin")
	go func() {
		closeCh := make(chan struct{})
		defer func() {
			close(closeCh)
			conn.Close()
		}()

		go func() {
			for {
				time.Sleep(time.Millisecond * 100)
				_ch, ok := s.parameterChan.Load(pluginName)
				if !ok {
					continue
				}
				ch := _ch.(chan parameterInfo)
				select {
				case <-closeCh:
					log.Println("Close")
					return
				case p := <-ch:
					if data, err := json.Marshal(wsMessage{
						Method:         "parameters",
						Plugin:         pluginName,
						ParameterName:  p.Name,
						ParameterValue: p.Value,
					}); err == nil {
						p, ok := config.C.Plugins[pluginName]
						if !ok {
							continue
						}
						if !p.AvatarBind(nowAvatar) {
							continue
						}
						_ = wsutil.WriteServerText(conn, data)
					}
				}
			}
		}()

		go func() {
			for {
				time.Sleep(time.Millisecond * 100)
				_ch, ok := s.avatarChangeChan.Load(pluginName)
				if !ok {
					continue
				}
				ch := _ch.(chan bool)
				select {
				case <-closeCh:
					log.Println("Close")
					return
				case change := <-ch:
					if change {
						if data, err := json.Marshal(wsMessage{
							Method: "avatar_change",
							Plugin: pluginName,
						}); err == nil {
							p, ok := config.C.Plugins[pluginName]
							if !ok {
								continue
							}
							if !p.AvatarBind(nowAvatar) {
								continue
							}
							_ = wsutil.WriteServerText(conn, data)
						}
					}
				}
			}
		}()

		for {
			msg, _, err := wsutil.ReadClientData(conn)
			if err != nil && !errors.Is(err, io.EOF) {
				log.Println(err)
				return
			}
			//log.Println("[WebSocket Message]", string(msg))
			var value wsMessage
			if json.Unmarshal(msg, &value) == nil {
				s.messageHandler(value, conn)
			}
		}
	}()
}

func (s *WSServer) messageHandler(msg wsMessage, conn net.Conn) {
	switch msg.Method {
	case "send":
		// 禁止未注册的插件发送消息
		_, ok := plugins[msg.Plugin]
		if !ok {
			return
		}

		valueStr, ok := msg.Value.(string)
		if ok && strings.Index(valueStr, "/avatar") == 0 {
			// 禁止向非绑定Avatar发送消息
			p, ok := config.C.Plugins[msg.Plugin]
			if !ok {
				return
			}
			if !p.AvatarBind(nowAvatar) {
				return
			}
		}

		m := osc.NewMessage(msg.Addr)
		switch v := msg.Value.(type) {
		case float32:
			m.Append(v)
		case float64:
			m.Append(float32(v))
		case int:
			m.Append(float32(v))
		case bool:
			m.Append(v)
		default:
			log.Println("[WebSocket Message]", "Unknown type", reflect.TypeOf(v))
			return
		}
		_ = s.osc.Send(m)

	case "get_options":
		if p, ok := config.C.Plugins[msg.Plugin]; ok {
			if data, err := json.Marshal(wsMessage{
				Method:  "get_options",
				Plugin:  msg.Plugin,
				Options: p.Options(),
			}); err == nil {
				_ = wsutil.WriteServerText(conn, data)
			}
		}

	case "listen_parameters":
		if p, ok := plugins[msg.Plugin]; ok {
			pluginsParameters.Store(p.Name, msg.Parameters)
			s.parameterChan.Store(p.Name, make(chan parameterInfo, 10))
		}

	case "listen_avatar_change":
		if p, ok := plugins[msg.Plugin]; ok {
			pluginsAvatarChange.Store(p.Name, false)
			s.avatarChangeChan.Store(p.Name, make(chan bool))
		}

	default:
		log.Println("[WebSocket Message]", "Unknown method", msg.Method)
	}
}

func (s *WSServer) Listen() error {
	logger.App("OSCServer").Printf("Listening on %s:%d\n", s.hostname, s.port)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.hostname, s.port), http.HandlerFunc(s.handle))
}
