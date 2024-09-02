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
	"github.com/mwantia/coredns-consulrpz-plugin/responses"
)

func HandlePoliciesParallel(state request.Request, ctx context.Context, _policies []policies.Policy) (*policies.Policy, *responses.PolicyResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	resultChannel := make(chan struct {
		Policy   *policies.Policy
		Response *responses.PolicyResponse
	}, len(_policies))
	errorChannel := make(chan error, len(_policies))

	for _, policy := range _policies {
		wg.Add(1)

		go func(pol policies.Policy) {
			defer wg.Done()

			start := time.Now()
			response, err := HandlePolicyResponse(state, ctx, policy)
			duration := time.Since(start).Seconds()

			metrics.MetricPolicyExecutionTime(policy.Name, duration)

			if err != nil {
				if !errors.Is(err, context.Canceled) {
					qname := dns.Fqdn(state.Name())
					logging.Log.Errorf("Unable to handle request for '%s': %s", qname, err)
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
					Response *responses.PolicyResponse
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
