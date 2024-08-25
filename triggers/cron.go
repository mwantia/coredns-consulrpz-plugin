package triggers

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/coredns/coredns/request"
	"github.com/robfig/cron/v3"
)

var (
	CronCompileCache = make(map[string]*cron.Schedule)
	CronCompileMutex sync.RWMutex
)

func MatchCronTrigger(state request.Request, value json.RawMessage) (bool, error) {
	now := time.Now()

	var expressions []string
	if err := json.Unmarshal(value, &expressions); err != nil {
		return false, err
	}

	for _, expression := range expressions {
		schedule, err := GetCachedCron(expression)
		if err != nil {
			return false, err
		}
		nextTime := schedule.Next(now)
		prevTime := schedule.Next(now.Add(-time.Minute))

		if now.After(prevTime) && now.Before(nextTime) {
			return true, nil
		}
	}

	return false, nil
}

func GetCachedCron(expression string) (cron.Schedule, error) {
	CronCompileMutex.RLock()
	schedule, exists := CronCompileCache[expression]
	CronCompileMutex.RUnlock()

	if exists {
		return *schedule, nil
	}

	CronCompileMutex.Lock()
	defer CronCompileMutex.Unlock()

	schedule, exists = CronCompileCache[expression]
	if exists {
		return *schedule, nil
	}

	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	parsedSchedule, err := parser.Parse(expression)
	if err != nil {
		return nil, err
	}

	CronCompileCache[expression] = &parsedSchedule
	return parsedSchedule, nil
}
