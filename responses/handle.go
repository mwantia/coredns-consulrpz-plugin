package responses

import (
	"context"
	"strings"

	"github.com/coredns/coredns/request"
	"github.com/mwantia/coredns-consulrpz-plugin/policies"
)

func HandleResponse(state request.Request, ctx context.Context, rule policies.PolicyRule) (*PolicyResponse, error) {
	response := &PolicyResponse{}
	for _, rresponse := range rule.Responses {
		alias := rresponse.GetAliasType()

		switch alias {
		case "deny":
			response.Deny = true

		case "fallthrough":
			response.Fallthrough = true

		case "code":
			if err := AppendRcodeToResponse(state, rresponse, response); err != nil {
				return response, err
			}

		case "record":
			if err := AppendRecordToResponse(state, rresponse, response); err != nil {
				return response, err
			}
		}

		if strings.HasPrefix(alias, "inaddr_") {
			if err := AppendInAddrToResponse(state, rresponse, response); err != nil {
				return response, err
			}
		}
	}

	return response, nil
}
