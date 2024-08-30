package rpz

import (
	"github.com/coredns/coredns/request"
	"github.com/mwantia/coredns-rpz-plugin/policies"
)

func HandleActionResponse(state request.Request, action policies.RuleAction, response *Response) (bool, error) {
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
