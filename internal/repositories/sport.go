package repositories

import (
	"betassist.ru/bookmaker-game-parser/internal/database"
	"betassist.ru/bookmaker-game-parser/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var Sport sportRepository

type sportRepository struct {
	master, slave   *gorm.DB
	AvailableSports map[string]uuid.UUID
}

func (r *sportRepository) Register() {
	Sport = sportRepository{
		master: database.GetMaster(),
		slave:  database.GetSlave(),
	}

	r.AvailableSports = r.loadAvailableSports()
}

func (r *sportRepository) loadAvailableSports() map[string]uuid.UUID {
	var sports []models.Sport

	r.slave.Model(new(models.Sport)).Select("id", "bookmaker_slug").Find(&sports)

	sportsMap := make(map[string]uuid.UUID, len(sports))

	for _, sport := range sports {
		sportsMap[sport.BookmakerSlug] = sport.ID
	}

	return sportsMap
}
