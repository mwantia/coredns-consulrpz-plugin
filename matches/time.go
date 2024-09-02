package matches

import (
	"context"
	"encoding/json"
	"time"

	"github.com/coredns/coredns/request"
)

type TimeData struct {
	TimeRanges []TimeRange
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

func ProcessTimeData(value json.RawMessage) (interface{}, error) {
	var timeranges []struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}
	if err := json.Unmarshal(value, &timeranges); err != nil {
		return nil, err
	}

	data := TimeData{}

	for _, timerange := range timeranges {
		start, err := time.Parse("15:04", timerange.Start)
		if err != nil {
			return nil, err
		}

		end, err := time.Parse("15:04", timerange.End)
		if err != nil {
			return false, err
		}

		data.TimeRanges = append(data.TimeRanges, TimeRange{
			Start: start,
			End:   end,
		})
	}
	return data, nil
}

func MatchTime(state request.Request, ctx context.Context, data TimeData) (bool, error) {
	now := time.Now()

	for _, timerange := range data.TimeRanges {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
			year, month, day := now.Date()

			start := time.Date(year, month, day, timerange.Start.Hour(), timerange.Start.Minute(), 0, 0, now.Location())
			end := time.Date(year, month, day, timerange.End.Hour(), timerange.End.Minute(), 0, 0, now.Location())

			if end.Before(start) {
				end = end.Add(24 * time.Hour)
			}

			if (now.After(start) || now.Equal(start)) && now.Before(end) {
				return true, nil
			}
		}
	}

	return false, nil
}
