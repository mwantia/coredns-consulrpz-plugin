package consul

import (
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

type ConsulWatchHandler func(api.KVPairs) error

func WatchConsulKVPrefix(address, token, prefix string, fn ConsulWatchHandler) error {
	def := CreateConsulConfig(address, token)
	params := map[string]interface{}{
		"type":   "keyprefix",
		"token":  def.Token,
		"prefix": prefix,
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

		fn(kv)
	}

	channel := make(chan error, 1)

	go func() {
		if err := watcher.Run(def.Address); err != nil {
			channel <- err
		}
		close(channel)
	}()

	select {
	case err := <-channel:
		return err
	case <-time.After(100 * time.Millisecond):
		return nil
	}
}
