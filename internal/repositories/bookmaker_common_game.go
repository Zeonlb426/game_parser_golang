package repositories

import (
	"betassist.ru/bookmaker-game-parser/internal/database"
	"betassist.ru/bookmaker-game-parser/internal/enums"
	"betassist.ru/bookmaker-game-parser/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var BookmakerCommonGame bookmakerCommonGameRepository

type bookmakerCommonGameRepository struct{ master, slave *gorm.DB }

func (bookmakerCommonGameRepository) Register() {
	BookmakerCommonGame = bookmakerCommonGameRepository{
		master: database.GetMaster(),
		slave:  database.GetSlave(),
	}
}

func (r *bookmakerCommonGameRepository) Upsert(id int, status enums.SearchStatus) models.BookmakerCommonGame {
	team := models.BookmakerCommonGame{
		ID:     id,
		Status: status,
	}

	r.master.Model(new(models.BookmakerCommonGame)).
		Clauses(clause.OnConflict{Columns: []clause.Column{
			{Name: "id"},
		}, DoNothing: true}).
		Create(&team)

	return team
}
