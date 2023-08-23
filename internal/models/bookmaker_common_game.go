package models

import (
	"betassist.ru/bookmaker-game-parser/internal/enums"
	"time"

	"betassist.ru/bookmaker-game-parser/internal/database/helpers"
)

type BookmakerCommonGame struct {
	ID     int                `gorm:"type:int(4);primaryKey;autoIncrement:false;"`
	Status enums.SearchStatus `gorm:"type:search_status"`

	CreatedAt time.Time `gorm:"not null;"`
	UpdatedAt time.Time `gorm:"not null;"`
}

func (m BookmakerCommonGame) TableName() string {
	return helpers.WithPrefix("bookmaker_common_games")
}
