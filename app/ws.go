package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
	"vrchat-osc-manager/internal/config"
	"vrchat-osc-manager/internal/logger"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/hypebeast/go-osc/osc"
)

type (
	WSServer struct {
		hostname string
		port     int
		osc      *OSC
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
		Avatar         string         `json:"avatar,omitempty"`
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

	plugin, ok := config.C.Plugins[pluginName]
	if !ok {
		return
	}

	ctx := context.Background()
	avatarChangeSub := hub.Sub(ctx, "avatar_change", func(msg *Message) {
		if plugin.ListenAvatarChange {
			avatar, ok := plugin.CheckAvatarBind(msg.avatar)
			if ok {
				if data, err := json.Marshal(wsMessage{
					Method: "avatar_change",
					Plugin: pluginName,
					Avatar: avatar,
				}); err == nil {
					_ = wsutil.WriteServerText(conn, data)
				}
			}
		}
	})
	defer hub.UnSub(avatarChangeSub)

	parametersSub := hub.Sub(ctx, "parameters", func(msg *Message) {
		if data, err := json.Marshal(wsMessage{
			Method:         "parameters",
			Plugin:         pluginName,
			ParameterName:  msg.parameter.name,
			ParameterValue: msg.parameter.value,
		}); err == nil {
			if _, ok = plugin.CheckAvatarBind(nowAvatar); ok {
				_ = wsutil.WriteServerText(conn, data)
			}
		}
	})
	defer hub.UnSub(parametersSub)

	closeCh := make(chan struct{})
	defer func() {
		close(closeCh)
		conn.Close()
	}()

	for {
		msg, _, err := wsutil.ReadClientData(conn)
		if err != nil && !errors.Is(err, io.EOF) {
			logger.App("WebSocket").Error(err)
			return
		}
		//log.Println("[WebSocket Message]", string(msg))
		var value wsMessage
		if json.Unmarshal(msg, &value) == nil {
			s.messageHandler(value, conn)
		}
	}
}

func (s *WSServer) messageHandler(msg wsMessage, conn net.Conn) {
	plugin, ok := config.C.Plugins[msg.Plugin]
	if !ok {
		return
	}
	switch msg.Method {
	case "send":
		valueStr, ok := msg.Value.(string)
		if ok && strings.Index(valueStr, "/avatar") == 0 {
			// 禁止向非绑定Avatar发送消息
			p, ok := config.C.Plugins[msg.Plugin]
			if !ok {
				break
			}
			if _, ok = p.CheckAvatarBind(nowAvatar); !ok {
				break
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
		if data, err := json.Marshal(wsMessage{
			Method:  "get_options",
			Plugin:  msg.Plugin,
			Options: plugin.Options(),
		}); err == nil {
			_ = wsutil.WriteServerText(conn, data)
		}

	case "listen_parameters":
		plugin.SetListenParameters(msg.Parameters)

	case "listen_avatar_change":
		plugin.ListenAvatarChange = true

	default:
		log.Println("[WebSocket Message]", "Unknown method", msg.Method)
	}
}

func (s *WSServer) Listen() error {
	logger.App("WebSocket").Printf("Listening on %s:%d\n", s.hostname, s.port)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.hostname, s.port), http.HandlerFunc(s.handle))
}
