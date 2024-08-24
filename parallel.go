package rpz

import (
	"context"
	"sync"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/mwantia/coredns-rpz-plugin/logging"
)

func (p RpzPlugin) HandlePoliciesParallel(state request.Request, ctx context.Context, request *dns.Msg) (*Policy, *Response, error) {
	var wg sync.WaitGroup

	resultChan := make(chan struct {
		Policy   *Policy
		Response *Response
	}, len(p.Config.Policies))
	errorChan := make(chan error, len(p.Config.Policies))

	for _, policy := range p.Config.Policies {
		wg.Add(1)
		go func(pol Policy) {
			defer wg.Done()
			response, err := HandlePolicy(state, ctx, request, policy)
			if err != nil {
				logging.Log.Errorf("Unable to handle request for '%s': %s", dns.Fqdn(state.Name()), err)

				errorChan <- err
				return
			}

			if response != nil {
				resultChan <- struct {
					Policy   *Policy
					Response *Response
				}{&pol, response}
			}
		}(policy)
	}

	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	case err := <-errorChan:
		return nil, nil, err
	case result := <-resultChan:
		return result.Policy, result.Response, nil
	}
}

func HandlePolicy(state request.Request, ctx context.Context, r *dns.Msg, policy Policy) (*Response, error) {
	logging.Log.Debugf("Handling policy named '%s'", policy.Name)

	for _, rule := range policy.Rules {
		if response, err := HandlePolicyRule(state, ctx, r, rule); response != nil || err != nil {
			return response, err
		}
	}

	return nil, nil
}

func HandlePolicyRule(state request.Request, ctx context.Context, r *dns.Msg, rule PolicyRule) (*Response, error) {
	for _, trigger := range rule.Triggers {
		if handled, err := HandleTrigger(state, trigger); !handled || err != nil {
			return nil, err
		}
	}

	if response, err := HandleResponse(state, ctx, r, rule); response != nil || err != nil {
		return response, err
	}

	return nil, nil
}
