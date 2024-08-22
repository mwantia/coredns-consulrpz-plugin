package rpz

import (
	"time"

	"github.com/hashicorp/consul/api"
)

type ConsulConfig struct {
	Prefix  string
	Address string
	Token   string
}

func CreateConsulClient(config *ConsulConfig) (*api.Client, error) {
	def := api.DefaultConfig()
	def.Address = config.Address
	def.Token = config.Token

	client, err := api.NewClient(def)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetConsulKVPairs(config ConsulConfig) (api.KVPairs, float64, error) {
	start := time.Now()
	options := &api.QueryOptions{
		UseCache:          true,
		MaxAge:            time.Minute,
		StaleIfError:      10 * time.Second,
		RequireConsistent: false,
		AllowStale:        true,
	}

	client, err := CreateConsulClient(&config)
	if err != nil {
		return nil, time.Since(start).Seconds(), err
	}

	pairs, _, err := client.KV().List(config.Prefix, options)
	if err != nil {
		return nil, time.Since(start).Seconds(), err
	}

	return pairs, time.Since(start).Seconds(), nil
}
