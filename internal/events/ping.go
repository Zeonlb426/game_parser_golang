package events

import (
	"errors"
)

const (
	queueCountLimit int = 1000
)

type Ping struct {
	Cmd             string  `json:"cmd"`
	LastDelay100avg float64 `json:"last_delay_100avg"`
	QueueCount      int     `json:"queue_count"`
}

type Pong struct {
	Cmd string `json:"cmd"`
}

func (e Ping) Validate() error {
	if e.QueueCount > queueCountLimit {
		return errors.New("reached queue count limit")
	}

	return nil
}
