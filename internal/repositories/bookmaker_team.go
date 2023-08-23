package repositories

import (
	"betassist.ru/bookmaker-game-parser/internal/database"
	"betassist.ru/bookmaker-game-parser/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var BookmakerTeam bookmakerTeamRepository

type bookmakerTeamRepository struct{ master, slave *gorm.DB }

func (bookmakerTeamRepository) Register() {
	BookmakerTeam = bookmakerTeamRepository{
		master: database.GetMaster(),
		slave:  database.GetSlave(),
	}
}

func (r *bookmakerTeamRepository) Upsert(bookmakerID uuid.UUID, name string) models.BookmakerTeam {
	team := models.BookmakerTeam{
		BookmakerID: bookmakerID,
		Name:        name,
	}

	r.master.Model(new(models.BookmakerTeam)).
		Clauses(clause.OnConflict{Columns: []clause.Column{
			{Name: "bookmaker_id"},
			{Name: "name"},
		}, DoNothing: true}).
		Create(&team)

	return team
}

func (r *bookmakerTeamRepository) Select(bookmakerID uuid.UUID, name string) (models.BookmakerTeam, error) {
	var bookmakerTeam models.BookmakerTeam

	tx := r.slave.Model(new(models.BookmakerTeam)).Select("id").Where("bookmaker_id = ?", bookmakerID).Where("name = ?", name).Take(&bookmakerTeam)

	return bookmakerTeam, tx.Error
}
