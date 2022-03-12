package app

import (
	"encoding/json"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/hypebeast/go-osc/osc"
	"log"
	"net/http"
	"reflect"
	"vrchat-osc-manager/internal/config"
)

type wsMessage struct {
	Method string      `json:"method"`
	Addr   string      `json:"addr"`
	Value  interface{} `json:"value"`
}

func wsServer() {
	_ = http.ListenAndServe(
		fmt.Sprintf("%s:%d", config.C.WebSocket.Hostname, config.C.WebSocket.Port),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, _, _, err := ws.UpgradeHTTP(r, w)
			if err != nil {
				log.Println(err)
				return
			}
			go func() {
				defer conn.Close()

				for {
					msg, _, err := wsutil.ReadClientData(conn)
					if err != nil {
						log.Println(err)
						return
					}
					//log.Println("[WebSocket Message]", string(msg))
					var value wsMessage
					if json.Unmarshal(msg, &value) == nil {
						messageHandler(value)
					}
				}
			}()
		}))
}

func messageHandler(msg wsMessage) {
	switch msg.Method {
	case "send":
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
		_ = oscClient.Send(m)
	default:
		log.Println("[WebSocket Message]", "Unknown method", msg.Method)
	}
}
