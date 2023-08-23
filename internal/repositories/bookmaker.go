package repositories

import (
	"betassist.ru/bookmaker-game-parser/internal/database"
	"betassist.ru/bookmaker-game-parser/internal/models"
	"gorm.io/gorm"
)

var Bookmaker bookmakerRepository

type bookmakerRepository struct{ master, slave *gorm.DB }

func (bookmakerRepository) Register() {
	Bookmaker = bookmakerRepository{
		master: database.GetMaster(),
		slave:  database.GetSlave(),
	}
}

func (r *bookmakerRepository) Select(slug string) (models.Bookmaker, error) {
	var bookmaker models.Bookmaker

	tx := r.slave.Model(new(models.Bookmaker)).Select("id").Where("slug = ?", slug).Take(&bookmaker)

	return bookmaker, tx.Error
}
