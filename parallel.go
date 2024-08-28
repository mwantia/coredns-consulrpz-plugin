package rpz

import (
	"context"
	"sync"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/mwantia/coredns-rpz-plugin/logging"
)

func (p RpzPlugin) HandlePoliciesParallel(state request.Request, ctx context.Context, request *dns.Msg) (*Policy, *Response, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	resultChannel := make(chan struct {
		Policy   *Policy
		Response *Response
	}, len(p.Config.Policies))
	errorChannel := make(chan error, len(p.Config.Policies))

	for _, policy := range p.Config.Policies {
		wg.Add(1)
		go func(pol Policy) {
			defer wg.Done()
			response, err := HandlePolicy(state, ctx, request, policy)
			if err != nil {
				logging.Log.Errorf("Unable to handle request for '%s': %s", dns.Fqdn(state.Name()), err)

				select {
				case errorChannel <- err:
				case <-ctx.Done():
				}
				return
			}

			if response != nil {
				select {
				case resultChannel <- struct {
					Policy   *Policy
					Response *Response
				}{&policy, response}:
				case <-ctx.Done():
				}
			}
		}(policy)
	}

	go func() {
		wg.Wait()

		close(resultChannel)
		close(errorChannel)
	}()

	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()

	case err := <-errorChannel:
		return nil, nil, err

	case result := <-resultChannel:
		return result.Policy, result.Response, nil
	}
}

func HandlePolicy(state request.Request, ctx context.Context, r *dns.Msg, policy Policy) (*Response, error) {
	logging.Log.Debugf("Handling policy named '%s'", policy.Name)

	for _, rule := range policy.Rules {
		if response, err := HandlePolicyRule(state, ctx, r, policy, rule); response != nil || err != nil {
			return response, err
		}
	}

	return nil, nil
}

func HandlePolicyRule(state request.Request, ctx context.Context, r *dns.Msg, policy Policy, rule PolicyRule) (*Response, error) {
	for _, trigger := range rule.Triggers {
		if handled, err := HandleTrigger(state, trigger); !handled || err != nil {
			return nil, err
		}

		alias := trigger.GetAliasType()
		MetricTriggerMatchCount(policy.Name, alias)
	}

	if response, err := HandleResponse(state, ctx, r, rule); response != nil || err != nil {
		return response, err
	}

	return nil, nil
}
