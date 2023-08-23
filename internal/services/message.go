package services

import (
	"betassist.ru/bookmaker-game-parser/internal/contracts"
	"betassist.ru/bookmaker-game-parser/internal/events"
	"betassist.ru/bookmaker-game-parser/internal/utils"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/valyala/fastjson"
	"strconv"
	"strings"
	"time"
)

const (
	defaultUpdateEventLength int = 5
)

var (
	pingName        = []byte("ping")
	updateEventName = []byte("update_event")
)

var jsonParser = fastjson.Parser{}

type Message struct {
	Name string          `json:"name"`
	Data contracts.Event `json:"data"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	v, errParse := jsonParser.ParseBytes(data)

	if nil != errParse {
		return errParse
	}

	if v.Exists("cmd") {
		if 0 == bytes.Compare(pingName, v.GetStringBytes("cmd")) {
			m.Name = string(pingName)

			m.Data = events.Ping{
				Cmd:             string(pingName),
				LastDelay100avg: v.GetFloat64("last_delay_100avg"),
				QueueCount:      v.GetInt("queue_count"),
			}

			return nil
		}

		return errors.New("unknown event")
	}

	array, errArray := v.Array()

	if nil != errArray {
		return errArray
	}

	if defaultUpdateEventLength <= len(array) {
		name := array[1].GetStringBytes()

		if 0 == bytes.Compare(updateEventName, name) {
			m.Name = string(updateEventName)

			updateEventData := events.UpdateEvent{
				Bookmaker: string(array[0].GetStringBytes()),
				ID:        string(array[2].GetStringBytes()),
			}

			if fastjson.TypeString == array[3].Type() {
				globalIDData := string(array[3].GetStringBytes())
				globalIDSlice := strings.Split(globalIDData, ",")

				if 2 == len(globalIDSlice) {
					globalID, _ := strconv.Atoi(globalIDSlice[0])
					updateEventData.GlobalID = &globalID

					teamSwappedInt, _ := strconv.Atoi(globalIDSlice[1])
					updateEventData.TeamSwapped = utils.IntToBool(teamSwappedInt)
				}
			}

			eventObject := array[4]

			updateEventData.Data = events.Event{
				ID:            string(eventObject.GetStringBytes("bk_event_id")),
				NativeID:      string(eventObject.GetStringBytes("bk_event_native_id")),
				BookmakerName: string(eventObject.GetStringBytes("bk_name")),
				SportID:       string(eventObject.GetStringBytes("sport_id")),
				Sport:         string(eventObject.GetStringBytes("sport")),
				EventName:     string(eventObject.GetStringBytes("event_name")),
				Team1:         string(eventObject.GetStringBytes("team1")),
				Team2:         string(eventObject.GetStringBytes("team2")),
				LeagueName:    string(eventObject.GetStringBytes("league_name")),
				DirectLink:    string(eventObject.GetStringBytes("direct_link")),
				Score:         string(eventObject.GetStringBytes("score")),
			}

			meta := eventObject.GetStringBytes("meta")
			updateEventData.Data.MetaRaw = meta
			errMeta := json.Unmarshal(meta, &updateEventData.Data.Meta)

			if nil != errMeta {
				return errMeta
			}

			addedAtStr := string(eventObject.GetStringBytes("added_at"))
			addedAtInt, err := strconv.ParseInt(addedAtStr, 10, 64)

			if nil == err && 0 != addedAtInt {
				addedAtTime := time.Unix(addedAtInt, 0)
				updateEventData.Data.AddedAt = &addedAtTime
			}

			m.Data = updateEventData

			return nil
		}
	}

	return errors.New("unknown event")
}
