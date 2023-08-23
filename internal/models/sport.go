package models

import (
	"betassist.ru/bookmaker-game-parser/internal/database/helpers"
	"github.com/google/uuid"
)

type Sport struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();"`

	Name           string `gorm:"type:varchar(255);not null;"`
	Slug           string `gorm:"type:varchar(255);not null;"`
	ChampionatSlug string `gorm:"type:varchar(255);not null;"`
	BookmakerSlug  string `gorm:"type:varchar(255);not null;"`
}

func (m Sport) TableName() string {
	return helpers.WithPrefix("sports")
}
