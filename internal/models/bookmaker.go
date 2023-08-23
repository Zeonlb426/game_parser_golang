package models

import (
	"betassist.ru/bookmaker-game-parser/internal/database/helpers"
	"github.com/google/uuid"
	"time"
)

type Bookmaker struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();"`

	Name string `gorm:"type:varchar(100);not null;"`
	Slug string `gorm:"type:varchar(100);uniqueIndex;not null;"`

	CreatedAt time.Time `gorm:"not null;"`
	UpdatedAt time.Time `gorm:"not null;"`
}

func (m Bookmaker) TableName() string {
	return helpers.WithPrefix("bookmakers")
}
