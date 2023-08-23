package events

import (
	"encoding/json"
	"errors"
	"time"

	"betassist.ru/bookmaker-game-parser/internal/utils"
)

const (
	SportFootball string = "soccer"
	SportHockey   string = "hockey"
)

var availableSports = []string{SportFootball, SportHockey}

type Meta struct {
	StartAt int64 `json:"start_at"`
}

type Event struct {
	ID            string          `json:"bk_event_id"`
	NativeID      string          `json:"bk_event_native_id"`
	BookmakerName string          `json:"bk_name"`
	SportID       string          `json:"sport_id"`
	Sport         string          `json:"sport"`
	EventName     string          `json:"event_name"`
	Team1         string          `json:"team1"`
	Team2         string          `json:"team2"`
	LeagueName    string          `json:"league_name"`
	DirectLink    string          `json:"direct_link"`
	AddedAt       *time.Time      `json:"added_at"`
	Meta          Meta            `json:"meta"`
	MetaRaw       json.RawMessage `json:"meta_raw"`
	Score         string          `json:"score"`
}

type UpdateEvent struct {
	Bookmaker   string `json:"bookmaker"`
	ID          string `json:"id"`
	GlobalID    *int   `json:"global_id"`
	TeamSwapped bool   `json:"team_swapped"`
	Data        Event  `json:"data"`
}

func (e UpdateEvent) Validate() error {
	if nil == e.GlobalID {
		return errors.New("empty global ID")
	}

	if !utils.StringInSlice(e.Data.Sport, availableSports) {
		return errors.New("unknown sport")
	}

	return nil
}
