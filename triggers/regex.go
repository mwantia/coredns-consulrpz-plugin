package triggers

import (
	"context"
	"encoding/json"
	"regexp"
	"sync"

	"github.com/coredns/coredns/request"
)

var (
	RegexCompileCache = make(map[string]*regexp.Regexp)
	RegexCompileMutex sync.RWMutex
)

func MatchRegexTrigger(state request.Request, ctx context.Context, value json.RawMessage) (bool, error) {
	var patterns []string
	if err := json.Unmarshal(value, &patterns); err != nil {
		return false, err
	}

	qname := state.Name()
	for _, pattern := range patterns {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
			if qname == pattern {
				return true, nil
			}

			regex, err := GetCachedRegex(pattern)
			if err != nil {
				return false, err
			}

			if regex.MatchString(qname) {
				return true, nil
			}
		}
	}

	return false, nil
}

func GetCachedRegex(pattern string) (*regexp.Regexp, error) {
	RegexCompileMutex.RLock()
	regex, exists := RegexCompileCache[pattern]
	RegexCompileMutex.RUnlock()

	if exists {
		return regex, nil
	}

	RegexCompileMutex.Lock()
	defer RegexCompileMutex.Unlock()

	regex, exists = RegexCompileCache[pattern]
	if exists {
		return regex, nil
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	RegexCompileCache[pattern] = regex
	return regex, nil
}
