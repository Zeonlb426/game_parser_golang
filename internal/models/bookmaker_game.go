package models

import (
	"betassist.ru/bookmaker-game-parser/internal/database/helpers"
	"github.com/google/uuid"
	"time"
)

type BookmakerGame struct {
	ID       string `gorm:"type:varchar(255);primaryKey;autoIncrement:false;"`
	NativeID string `gorm:"type:varchar(255);not null;uniqueIndex:idx-bookmaker_games-bookmaker_id-native_id,priority:2;"`

	CommonGameID int                 `gorm:"type:int;not null;"`
	CommonGame   BookmakerCommonGame `gorm:"foreignKey:CommonGameID;"`

	RawName string `gorm:"type:varchar(255);not null;"`

	OwnersID   uuid.UUID     `gorm:"type:uuid;not null;"`
	OwnersTeam BookmakerTeam `gorm:"foreignKey:OwnersID;"`

	GuestsID  uuid.UUID     `gorm:"type:uuid;not null;"`
	GuestTeam BookmakerTeam `gorm:"foreignKey:GuestsID;"`

	TournamentID uuid.UUID           `gorm:"type:uuid;not null;"`
	Tournament   BookmakerTournament `gorm:"foreignKey:TournamentID;"`

	Link string `gorm:"type:varchar(255);not null;"`

	Meta []byte `gorm:"type:jsonb"`

	Score string `gorm:"type:varchar(255);not null;"`

	Datetime time.Time

	BookmakerID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx-bookmaker_games-bookmaker_id-native_id,priority:1;"`
	Bookmaker   Bookmaker `gorm:"foreignKey:BookmakerID;"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Swap bool
}

func (m BookmakerGame) TableName() string {
	return helpers.WithPrefix("bookmaker_games")
}
