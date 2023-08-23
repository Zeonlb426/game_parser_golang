package services

import (
	core "betassist.ru/bookmaker-game-parser/internal"
	"betassist.ru/bookmaker-game-parser/internal/enums"
	"betassist.ru/bookmaker-game-parser/internal/events"
	"betassist.ru/bookmaker-game-parser/internal/repositories"
	"github.com/rs/zerolog/log"
)

var Listener listenerService

type listenerService struct {
	updateEvent chan events.UpdateEvent
	ping        chan events.Ping
	close       chan struct{}
}

func init() {
	Listener.updateEvent = make(chan events.UpdateEvent, 30)
	Listener.ping = make(chan events.Ping, 1)
	Listener.close = make(chan struct{}, 1)
}

func (s *listenerService) Close() {
	s.close <- struct{}{}
}

func (s *listenerService) Listen() {
	for {
		select {
		case <-s.close:
			break
		case updateEvent := <-s.updateEvent:
			s.processEvent(updateEvent)
		case ping := <-s.ping:
			s.processPing(ping)
		}
	}
}

func (s *listenerService) processPing(ping events.Ping) {
	if ping.QueueCount < 100 {
		log.Info().Str("tag", core.WebsocketLogTag).Int("count", ping.QueueCount).Msg("message queue")
	} else if ping.QueueCount < 500 {
		log.Warn().Str("tag", core.WebsocketLogTag).Int("count", ping.QueueCount).Msg("message queue")
	} else {
		log.Error().Str("tag", core.WebsocketLogTag).Int("count", ping.QueueCount).Msg("message queue")
	}
}

func (s *listenerService) processEvent(event events.UpdateEvent) {
	sportID, ok := repositories.Sport.AvailableSports[event.Data.Sport]

	if !ok {
		return
	}

	team1 := Game.GetTeam(event.Data.Team1)
	team2 := Game.GetTeam(event.Data.Team2)
	tournament := Game.GetTournament(event.Data.LeagueName, sportID)
	commonGame := repositories.BookmakerCommonGame.Upsert(*event.GlobalID, enums.Pending)

	Game.StoreGame(team1, team2, tournament, commonGame, event)
}
