package consulrpz

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/hashicorp/consul/api"
	"github.com/mwantia/coredns-consulrpz-plugin/logging"
	"github.com/mwantia/coredns-consulrpz-plugin/metrics"
	"github.com/mwantia/coredns-consulrpz-plugin/policies"
)

type ConsulRpzPlugin struct {
	Next     plugin.Handler
	Cfg      *ConsulRpzConfig
	Consul   *api.Client
	Policies []policies.Policy
}

func init() {
	plugin.Register("consulrpz", setup)
}

func setup(c *caddy.Controller) error {
	c.OnStartup(func() error {
		return metrics.Register()
	})

	plug, err := CreatePlugin(c)
	if err != nil {
		logging.Log.Errorf("%v", err)
		return plugin.Error("consulrpz", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		plug.Next = next
		return plug
	})

	return nil
}
