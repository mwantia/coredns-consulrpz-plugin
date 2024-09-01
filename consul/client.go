package consul

import (
	"os"
	"time"

	"github.com/hashicorp/consul/api"
)

func CreateConsulConfig(address, token string) *api.Config {
	def := api.DefaultConfig()
	def.Address = address
	def.Token = token

	UpdateEnvConfig(def)
	return def
}

func CreateConsulClient(address, token string) (*api.Client, error) {
	def := CreateConsulConfig(address, token)

	client, err := api.NewClient(def)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetConsulKVPairs(client *api.Client, prefix string) (api.KVPairs, float64, error) {
	start := time.Now()
	options := &api.QueryOptions{
		UseCache:          true,
		MaxAge:            time.Minute,
		StaleIfError:      10 * time.Second,
		RequireConsistent: false,
		AllowStale:        true,
	}

	pairs, _, err := client.KV().List(prefix, options)
	duration := time.Since(start).Seconds()

	return pairs, duration, err
}

func UpdateEnvConfig(cfg *api.Config) {
	address := os.Getenv("CONSUL_HTTP_ADDR")
	if address != "" {
		cfg.Address = address
	}

	token := os.Getenv("CONSUL_HTTP_TOKEN")
	if address != "" {
		cfg.Token = token
	}
}
