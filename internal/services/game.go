package services

import (
	core "betassist.ru/bookmaker-game-parser/internal"
	"betassist.ru/bookmaker-game-parser/internal/events"
	"betassist.ru/bookmaker-game-parser/internal/models"
	"betassist.ru/bookmaker-game-parser/internal/repositories"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type gameService struct{}

var Game gameService

func (s *gameService) GetTeam(name string) models.BookmakerTeam {
	team := repositories.BookmakerTeam.Upsert(core.Config.App.BookmakerID, name)

	if uuid.Nil == team.ID {
		team, _ = repositories.BookmakerTeam.Select(core.Config.App.BookmakerID, name)
	}

	return team
}

func (s *gameService) GetTournament(name string, sport uuid.UUID) models.BookmakerTournament {
	tournament := repositories.BookmakerTournament.Upsert(core.Config.App.BookmakerID, name, sport)

	if uuid.Nil == tournament.ID {
		tournament, _ = repositories.BookmakerTournament.Select(core.Config.App.BookmakerID, name)
	}

	return tournament
}

func (s *gameService) StoreGame(team1 models.BookmakerTeam, team2 models.BookmakerTeam, tournament models.BookmakerTournament, commonGame models.BookmakerCommonGame, event events.UpdateEvent) {
	bookmakerGame := models.BookmakerGame{
		ID:           event.Data.ID,
		BookmakerID:  core.Config.App.BookmakerID,
		NativeID:     event.Data.NativeID,
		CommonGameID: commonGame.ID,
		RawName:      event.Data.EventName,
		TournamentID: tournament.ID,
		Link:         event.Data.DirectLink,
		Score:        event.Data.Score,
		Datetime:     time.Unix(event.Data.Meta.StartAt, 0),
		Swap:         event.TeamSwapped,
		OwnersID:     team1.ID,
		GuestsID:     team2.ID,
	}

	meta, err := json.Marshal(event.Data.Meta)

	if err != nil {
		meta = []byte("{}")
	}

	bookmakerGame.Meta = meta

	repositories.BookmakerGame.Upsert(bookmakerGame)
}
