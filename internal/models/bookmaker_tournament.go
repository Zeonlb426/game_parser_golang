package models

import (
	"betassist.ru/bookmaker-game-parser/internal/database/helpers"
	"github.com/google/uuid"
)

type BookmakerTournament struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();"`

	BookmakerID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:bookmaker_teams_bookmaker_id_name_unique,priority:1;"`
	Bookmaker   Bookmaker `gorm:"foreignKey:BookmakerID;"`

	SportID uuid.UUID `gorm:"type:uuid;not null;"`
	Sport   Sport     `gorm:"foreignKey:SportID;"`

	Name string `gorm:"type:varchar(255);not null;uniqueIndex:bookmaker_teams_bookmaker_id_name_unique,priority:2;"`
}

func (m BookmakerTournament) TableName() string {
	return helpers.WithPrefix("bookmaker_tournaments")
}
