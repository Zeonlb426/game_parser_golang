package repositories

import (
	"betassist.ru/bookmaker-game-parser/internal/database"
	"betassist.ru/bookmaker-game-parser/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var BookmakerTournament bookmakerTournamentRepository

type bookmakerTournamentRepository struct{ master, slave *gorm.DB }

func (bookmakerTournamentRepository) Register() {
	BookmakerTournament = bookmakerTournamentRepository{
		master: database.GetMaster(),
		slave:  database.GetSlave(),
	}
}

func (r *bookmakerTournamentRepository) Upsert(bookmakerID uuid.UUID, name string, sport uuid.UUID) models.BookmakerTournament {
	tournament := models.BookmakerTournament{
		BookmakerID: bookmakerID,
		Name:        name,
		SportID:     sport,
	}

	r.master.Model(new(models.BookmakerTournament)).
		Clauses(clause.OnConflict{Columns: []clause.Column{
			{Name: "bookmaker_id"},
			{Name: "name"},
		}, DoNothing: true}).
		Create(&tournament)

	return tournament
}

func (r *bookmakerTournamentRepository) Select(bookmakerID uuid.UUID, name string) (models.BookmakerTournament, error) {
	var bookmakerTournament models.BookmakerTournament

	tx := r.slave.Model(new(models.BookmakerTournament)).Select("id").Where("bookmaker_id = ?", bookmakerID).Where("name = ?", name).Take(&bookmakerTournament)

	return bookmakerTournament, tx.Error
}
