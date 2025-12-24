package ws

import (
	"context"
	"log"
	"net/http"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"

	"pulse_agent/internal/config"
	"pulse_agent/internal/terminal"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func ConnectAgentWS(ctx context.Context, cfg *config.Config, serverUUID string) error {
	header := http.Header{}
	header.Set("x-api-key", cfg.APIKey)

	wsURL := cfg.BackendURL
	wsURL = "ws" + wsURL[4:] + "/ws/agent"

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, header)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Println("Connected to backend WS")

	// 🔐 Register agent
	if err := conn.WriteJSON(Message{
		Type: "agent:register",
		Data: map[string]string{
			"server_uuid": serverUUID,
		},
	}); err != nil {
		return err
	}
	print(" registered the serverr")

	session, err := terminal.StartShell()
	print(" shell is started ", session)

	if err != nil {
		return err
	}
	defer session.Close()

	// 📤 Shell stdout → backend
	go session.ReadLoop(func(data []byte) {
		_ = conn.WriteJSON(Message{
			Type: "terminal:stdout",
			Data: string(data),
		})
	})

	// 📥 Backend → shell stdin
	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			return err
		}

		switch msg.Type {

		case "terminal:stdin":
			text, ok := msg.Data.(string)
			if ok {
				_ = session.Write([]byte(text))
			}

		case "terminal:resize":
			size, ok := msg.Data.(map[string]interface{})
			if !ok {
				continue
			}

			rows, rOk := size["rows"].(float64)
			cols, cOk := size["cols"].(float64)
			if !rOk || !cOk {
				continue
			}

			_ = pty.Setsize(session.Pty, &pty.Winsize{
				Rows: uint16(rows),
				Cols: uint16(cols),
			})
		}
	}
}
