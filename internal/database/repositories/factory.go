package repositories

import (
	"betassist.ru/bookmaker-game-parser/internal/repositories"
	"github.com/rs/zerolog/log"
)

func Register() {
	repositories.Bookmaker.Register()
	repositories.BookmakerTeam.Register()
	repositories.BookmakerTournament.Register()
	repositories.BookmakerCommonGame.Register()
	repositories.BookmakerGame.Register()
	repositories.Sport.Register()

	log.Debug().Msg("repositories registered")
}
