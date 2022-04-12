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
	"vrchat-osc-manager/internal/config"
	"vrchat-osc-manager/internal/utils"

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
		Method  string                 `json:"method"`
		Plugin  string                 `json:"plugin"`
		Addr    string                 `json:"addr,omitempty"`
		Value   interface{}            `json:"value,omitempty"`
		Options map[string]interface{} `json:"options,omitempty"`
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
	go func() {
		defer conn.Close()

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
		p, ok := config.C.Plugins[msg.Plugin]
		if !ok {
			return
		}
		if !p.AvatarBind(nowAvatar) {
			return
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

	default:
		log.Println("[WebSocket Message]", "Unknown method", msg.Method)
	}
}

func (s *WSServer) Listen() error {
	fmt.Printf("%s Listening on %s:%d\n", utils.Color(34, "[WebSocket]"), s.hostname, s.port)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.hostname, s.port), http.HandlerFunc(s.handle))
}
