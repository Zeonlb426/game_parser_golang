package repositories

import (
	"betassist.ru/bookmaker-game-parser/internal/database"
	"betassist.ru/bookmaker-game-parser/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var BookmakerGame bookmakerGameRepository

type bookmakerGameRepository struct{ master, slave *gorm.DB }

func (bookmakerGameRepository) Register() {
	BookmakerGame = bookmakerGameRepository{
		master: database.GetMaster(),
		slave:  database.GetSlave(),
	}
}

func (r *bookmakerGameRepository) Upsert(team models.BookmakerGame) models.BookmakerGame {
	r.master.Model(new(models.BookmakerGame)).
		Clauses(clause.OnConflict{Columns: []clause.Column{
			{Name: "bookmaker_id"},
			{Name: "native_id"},
		}, DoNothing: true}).
		Create(&team)

	return team
}
