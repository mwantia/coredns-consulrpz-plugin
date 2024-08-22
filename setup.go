package rpz

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	plugin.Register("rpz", setup)
}

func setup(c *caddy.Controller) error {
	c.OnStartup(func() error {
		prometheus.MustRegister(metricsRpzRequestDurationSeconds)
		prometheus.MustRegister(metricsQueryRequestsTotal)
		return nil
	})

	plug, err := CreatePlugin(c)
	if err != nil {
		return plugin.Error("rpz", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		plug.Next = next
		return plug
	})

	return nil
}
