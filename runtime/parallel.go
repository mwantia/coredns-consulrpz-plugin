package runtime

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/mwantia/coredns-consulrpz-plugin/logging"
	"github.com/mwantia/coredns-consulrpz-plugin/metrics"
	"github.com/mwantia/coredns-consulrpz-plugin/policies"
)

func HandlePoliciesParallel(state request.Request, ctx context.Context, request *dns.Msg, p []policies.Policy) (*policies.Policy, *Response, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	resultChannel := make(chan struct {
		Policy   *policies.Policy
		Response *Response
	}, len(p))
	errorChannel := make(chan error, len(p))

	for _, policy := range p {
		wg.Add(1)
		go func(pol policies.Policy) {
			defer wg.Done()
			start := time.Now()
			response, err := HandlePolicy(state, ctx, request, policy)
			duration := time.Since(start).Seconds()

			metrics.MetricPolicyExecutionTime(policy.Name, duration)

			if err != nil {
				if !errors.Is(err, context.Canceled) {
					logging.Log.Errorf("Unable to handle request for '%s': %s", dns.Fqdn(state.Name()), err)
				}

				select {
				case errorChannel <- err:
				case <-ctx.Done():
				}
				return
			}

			if response != nil {
				select {
				case resultChannel <- struct {
					Policy   *policies.Policy
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

func HandlePolicy(state request.Request, ctx context.Context, r *dns.Msg, policy policies.Policy) (*Response, error) {
	logging.Log.Debugf("Handling policy named '%s'", policy.Name)

	for _, rule := range policy.Rules {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if response, err := HandlePolicyRule(state, ctx, r, policy, rule); response != nil || err != nil {
				return response, err
			}
		}
	}

	return nil, nil
}

func HandlePolicyRule(state request.Request, ctx context.Context, r *dns.Msg, policy policies.Policy, rule policies.PolicyRule) (*Response, error) {
	for _, trigger := range rule.Triggers {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if handled, err := HandleTrigger(state, ctx, trigger); !handled || err != nil {
				return nil, err
			}

			alias := trigger.GetAliasType()
			metrics.MetricTriggerMatchCount(policy.Name, alias)
		}
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if response, err := HandleResponse(state, ctx, r, rule); response != nil || err != nil {
			return response, err
		}
	}

	return nil, nil
}
