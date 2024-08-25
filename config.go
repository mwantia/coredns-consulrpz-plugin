package rpz

import (
	"bytes"
	"os"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/plugin"
	"github.com/hashicorp/consul/api"
	"github.com/mwantia/coredns-rpz-plugin/logging"
)

type RpzPlugin struct {
	Next   plugin.Handler
	Config *RpzConfig
}

type RpzConfig struct {
	Policies []Policy
}

func CreatePlugin(c *caddy.Controller) (*RpzPlugin, error) {
	plug := &RpzPlugin{}

	config, err := CreateConfig(c)
	if err != nil {
		return nil, err
	}

	plug.Config = config

	return plug, nil
}

func CreateConfig(c *caddy.Controller) (*RpzConfig, error) {
	config := &RpzConfig{}

	n := 0
	for c.Next() {
		if n > 0 {
			return nil, c.Err("Unable to load config")
		}
		n++

		args := c.RemainingArgs()
		if len(args) >= 1 {
			logging.Log.Debugf("Available args: %v", args)
		}

		for c.NextBlock() {
			val := c.Val()
			args = c.RemainingArgs()

			if len(args) < 1 {
				return nil, c.Errf("config '%s' can't be empty", val)
			}

			switch val {
			case "consul":
				consul := ConsulConfig{
					Prefix:  args[0],
					Address: "http://127.0.0.1:8500",
					Token:   "",
				}

				if len(args) > 1 {
					consul.Address = args[1]
				}
				if len(args) > 2 {
					consul.Token = args[2]
				}

				pairs, _, err := GetConsulKVPairs(consul)
				if err != nil {
					logging.Log.Warningf("Unable to load consul prefix '%s': %v", args[0], err)
					continue
				}

				if err := config.ParseConsulKVPairs(pairs); err != nil {
					logging.Log.Warningf("Unable to parse consul kvpairs '%s': %v", args[0], err)
				}

				if err := WatchConsulPrefix(consul, func(pairs api.KVPairs) error {
					return config.ParseConsulKVPairs(pairs)
				}); err != nil {
					logging.Log.Warningf("Unable to load consul prefix '%s': %v", args[0], err)
				}

			case "policy":
				for _, a := range args {
					file, err := os.Open(a)
					if err != nil {
						logging.Log.Warningf("Unable to load file '%s': %v", a, err)
						continue
					}

					policy, err := ParsePolicyFile(file)
					if err != nil {
						logging.Log.Warningf("Unable to parse policy '%s': %v", a, err)
						continue
					}

					if err := config.UpdateNamedPolicies(policy); err != nil {
						logging.Log.Warningf("Unable to update policies: %v", err)
						continue
					}
				}
			}
		}
	}

	SortPolicies(config.Policies)
	return config, nil
}

func (c *RpzConfig) ParseConsulKVPairs(pairs api.KVPairs) error {
	for _, kv := range pairs {
		if err := c.ParseConsulKVPair(kv); err != nil {
			logging.Log.Warningf("Unable to parse kvpair '%s': %v", kv.Key, err)
			continue
		}
	}
	return nil
}

func (c *RpzConfig) ParseConsulKVPair(kv *api.KVPair) error {
	logging.Log.Infof("Parsing consul kvpair '%s'", kv.Key)
	reader := bytes.NewReader(kv.Value)

	policy, err := ParsePolicyFile(reader)
	if err != nil {
		return err
	}

	if err := c.UpdateNamedPolicies(policy); err != nil {
		return err
	}

	return nil
}

func (c *RpzConfig) UpdateNamedPolicies(policy *Policy) error {
	if policy != nil {
		for i := range c.Policies {
			if c.Policies[i].Name == policy.Name {
				logging.Log.Debugf("Checking hash for policy '%s':", policy.Name)
				logging.Log.Debugf("  Hash1: %s", c.Policies[i].Hash)
				logging.Log.Debugf("  Hash2: %s", policy.Hash)

				if c.Policies[i].Hash != policy.Hash {
					c.Policies[i].Priority = policy.Priority
					c.Policies[i].Rules = policy.Rules
					c.Policies[i].Hash = policy.Hash

					logging.Log.Debugf("Policy '%s' updated", policy.Name)
				}

				return nil
			}
		}

		logging.Log.Infof("Policy '%s' added to the list with hash ['%s']", policy.Name, policy.Hash)
		c.Policies = append(c.Policies, *policy)
	}
	return nil
}
