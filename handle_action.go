package rpz

import (
	"strings"

	"github.com/coredns/coredns/request"
)

var ActionTypeAliasMap = map[string]string{
	"rcode": "code",
}

func HandleActionResponse(state request.Request, action RuleAction, response *Response) (bool, error) {
	alias := action.GetAliasType()

	switch alias {
	case "deny":
		response.Deny = true
		return true, nil

	case "fallthrough":
		response.Fallthrough = true
		return true, nil

	case "code":
		if err := response.AppendRcode(state, action); err != nil {
			return false, err
		}

	case "record":
		if err := response.AppendRecord(state, action); err != nil {
			return false, err
		}
	}

	return false, nil
}

func (action RuleAction) GetAliasType() string {
	t := strings.ToLower(action.Type)
	if alias, exist := ActionTypeAliasMap[t]; exist {
		return alias
	}

	return t
}
