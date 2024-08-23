package triggers

import (
	"encoding/json"
	"regexp"
	"sync"

	"github.com/coredns/coredns/request"
)

var (
	RegexCompileCache = make(map[string]*regexp.Regexp)
	RegexCompileMutex sync.RWMutex
)

func MatchQNameTrigger(state request.Request, value json.RawMessage) (bool, error) {
	var domains []string
	if err := json.Unmarshal(value, &domains); err != nil {
		return false, err
	}

	qname := state.Name()
	for _, d := range domains {
		if qname == d {
			return true, nil
		}

		regex, err := GetCachedRegex(d)
		if err != nil {
			return false, err
		}

		if regex.MatchString(qname) {
			return true, nil
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
