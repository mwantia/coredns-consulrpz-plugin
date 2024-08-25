package triggers

import (
	"encoding/json"
	"time"

	"github.com/coredns/coredns/request"
)

type TimeTrigger struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func MatchTimeTrigger(state request.Request, value json.RawMessage) (bool, error) {
	now := time.Now()

	var tts []TimeTrigger
	if err := json.Unmarshal(value, &tts); err != nil {
		return false, err
	}

	for _, tt := range tts {
		start, err := time.Parse("15:04", tt.Start)
		if err != nil {
			return false, err
		}

		end, err := time.Parse("15:04", tt.End)
		if err != nil {
			return false, err
		}

		year, month, day := now.Date()
		start = time.Date(year, month, day, start.Hour(), start.Minute(), 0, 0, now.Location())
		end = time.Date(year, month, day, end.Hour(), end.Minute(), 0, 0, now.Location())

		if end.Before(start) {
			end = end.Add(24 * time.Hour)
		}

		if (now.After(start) || now.Equal(start)) && now.Before(end) {
			return true, nil
		}
	}

	return false, nil
}
