package responses

import (
	"context"
	"strings"

	"github.com/coredns/coredns/request"
	"github.com/mwantia/coredns-consulrpz-plugin/policies"
)

func HandleResponse(state request.Request, ctx context.Context, rule policies.PolicyRule) (*PolicyResponse, error) {
	response := &PolicyResponse{}
	for _, rr := range rule.Responses {
		alias := rr.GetAliasType()

		switch alias {
		case "deny":
			response.Deny = true

		case "fallthrough":
			response.Fallthrough = true

		case "extra":
			if err := AppendExtraToResponse(state, rr.Value, response); err != nil {
				return response, err
			}

		case "code":
			if err := AppendRcodeToResponse(state, rr.Value, response); err != nil {
				return response, err
			}

		case "record":
			if err := AppendRecordToResponse(state, rr.Value, response); err != nil {
				return response, err
			}
		}

		if strings.HasPrefix(alias, "inaddr_") {
			if err := AppendInAddrToResponse(state, alias, response); err != nil {
				return response, err
			}
		}
	}

	return response, nil
}
