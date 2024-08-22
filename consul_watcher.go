package rpz

import (
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/mwantia/coredns-rpz-plugin/logging"
)

type handler func(api.KVPairs) error

func WatchConsulPrefix(config ConsulConfig, fn handler) error {
	params := map[string]interface{}{
		"type":   "keyprefix",
		"prefix": config.Prefix,
		"token":  config.Token,
	}

	watcher, err := watch.Parse(params)
	if err != nil {
		return err
	}

	watcher.Handler = func(idx uint64, raw interface{}) {

		if raw == nil {
			return
		}

		kv, ok := raw.(api.KVPairs)
		if !ok || kv == nil {
			return
		}

		logging.Log.Debugf("Detected changes in Consul prefix '%s'", config.Prefix)
		fn(kv)
	}

	go func() {
		if err := watcher.Run(config.Address); err != nil {
			logging.Log.Errorf("Error running watch plan: %v", err)
		}
	}()

	return nil
}
