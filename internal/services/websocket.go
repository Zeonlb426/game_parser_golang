package services

import (
	core "betassist.ru/bookmaker-game-parser/internal"
	"betassist.ru/bookmaker-game-parser/internal/events"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/url"
)

const (
	cmdName string = "subscribe"
)

var Websocket websocketService

type websocketService struct {
	url  url.URL
	conn *websocket.Conn
}

type subscribeMessage struct {
	Cmd           string   `json:"cmd"`
	AuthKey       string   `json:"auth_key"`
	NeededBk      []string `json:"needed_bk"`
	SendEventsIds bool     `json:"send_events_ids"`
}

func init() {
	Websocket.url = url.URL{
		Scheme: core.Config.Websocket.Protocol,
		Host:   core.Config.Websocket.Host,
		Path:   core.Config.Websocket.Path,
	}
}

func (s *websocketService) Register() {
	var err error

	s.conn, _, err = websocket.DefaultDialer.Dial(s.url.String(), nil)

	if nil != err {
		panic(err)
	}
}

func (s *websocketService) sendInitialMessage() error {
	return s.conn.WriteJSON(subscribeMessage{
		Cmd:           cmdName,
		AuthKey:       core.Config.Websocket.AuthKey,
		NeededBk:      []string{fmt.Sprintf("%s:prematch", core.Config.Websocket.Bookmaker)},
		SendEventsIds: true,
	})
}

func (s *websocketService) Close() error {
	return s.conn.Close()
}

func (s *websocketService) Listen() {
	_ = s.sendInitialMessage()

	for {
		messageType, message, err := s.conn.ReadMessage()

		if nil != err {
			_ = s.Close()

			panic(err)
		}

		switch messageType {
		case websocket.CloseMessage:
			_ = s.Close()

			panic("close message from server")
		case websocket.PingMessage:
			// AM: no reaction is not correct: https://datatracker.ietf.org/doc/html/rfc6455
		case websocket.PongMessage:
			// AM: no reaction is not correct: https://datatracker.ietf.org/doc/html/rfc6455
		case websocket.TextMessage:
			s.processMessage(message)
		}
	}
}

func (s *websocketService) processMessage(data []byte) {
	if core.Config.Websocket.Debug {
		log.Debug().Str("tag", core.WebsocketLogTag).Bytes("data", data).Msg("new event")
	}

	var message Message

	if err := json.Unmarshal(data, &message); nil != err {
		return
	}

	if err := message.Data.Validate(); nil != err {
		return
	}

	switch message.Data.(type) {
	case events.Ping:
		Listener.ping <- message.Data.(events.Ping)
	case events.UpdateEvent:
		Listener.updateEvent <- message.Data.(events.UpdateEvent)
	default:
	}
}
