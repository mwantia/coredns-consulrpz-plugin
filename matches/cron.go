package matches

import (
	"context"
	"encoding/json"
	"time"

	"github.com/coredns/coredns/request"
	"github.com/robfig/cron/v3"
)

type CronData struct {
	Entries []struct {
		Expression string
		Schedule   cron.Schedule
	}
}

func ProcessCronData(value json.RawMessage) (interface{}, error) {
	var expressions []string
	if err := json.Unmarshal(value, &expressions); err != nil {
		return nil, err
	}

	data := CronData{}

	for _, expression := range expressions {
		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		schedule, err := parser.Parse(expression)
		if err != nil {
			return nil, err
		}

		data.Entries = append(data.Entries, struct {
			Expression string
			Schedule   cron.Schedule
		}{Expression: expression, Schedule: schedule})
	}

	return data, nil
}

func MatchCron(state request.Request, ctx context.Context, data CronData) (*MatchResult, error) {
	now := time.Now()

	for _, entry := range data.Entries {
		nextTime := entry.Schedule.Next(now)
		prevTime := entry.Schedule.Next(now.Add(-time.Minute))

		if now.After(prevTime) && now.Before(nextTime) {
			return &MatchResult{
				Handled: true,
				Data:    entry.Expression,
			}, nil
		}
	}

	return nil, nil
}
