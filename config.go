package consulrpz

import (
	"bytes"
	"fmt"

	"github.com/coredns/caddy"
	"github.com/hashicorp/consul/api"
	"github.com/mwantia/coredns-consulrpz-plugin/consul"
	"github.com/mwantia/coredns-consulrpz-plugin/logging"
	"github.com/mwantia/coredns-consulrpz-plugin/policies"
	"github.com/rschone/corefile2struct/pkg/corefile"
)

type ConsulRpzConfig struct {
	Arguments []string // This will store the prefix (allows for multiple entries)

	Address   string `cf:"address" default:"http://127.0.0.1:8500"`
	Token     string `cf:"token"`
	Watch     bool   `cf:"watch" default:"true"`
	Execution string `cf:"execution" default:"parallel" check:"oneOf(sequence|parallel)"`
}

func CreatePlugin(c *caddy.Controller) (*ConsulRpzPlugin, error) {
	p := &ConsulRpzPlugin{}

	var err error
	var cfg ConsulRpzConfig
	if err = corefile.Parse(c, &cfg); err != nil {
		return nil, err
	}

	p.Cfg = &cfg
	if p.Consul, err = consul.CreateConsulClient(cfg.Address, cfg.Token); err != nil {
		return nil, err
	}

	if err := LoadConsulKVPrefix(p); err != nil {
		return nil, err
	}
	policies.SortPolicies(p.Policies)

	return p, nil
}

func LoadConsulKVPrefix(plug *ConsulRpzPlugin) error {
	logging.Log.Infof("Load prefixes from arguments: %s", plug.Cfg.Arguments)
	// Each arguments passed onto consulrpz will be handled as prefix
	for _, prefix := range plug.Cfg.Arguments {
		pairs, _, err := consul.GetConsulKVPairs(plug.Consul, prefix)
		if err != nil {
			return err
		}

		if err := ParseConsulKVPairs(plug, pairs); err != nil {
			return err
		}

		if plug.Cfg.Watch {
			logging.Log.Infof("Watching prefix '%s' for new changes/updates", prefix)
			// Only watch if the config has been set to true
			if err := consul.WatchConsulKVPrefix(plug.Cfg.Address, plug.Cfg.Token, prefix, func(watchpairs api.KVPairs) error {
				logging.Log.Debugf("New update for prefix '%s'", prefix)
				return ParseConsulKVPairs(plug, watchpairs)
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func ParseConsulKVPairs(plug *ConsulRpzPlugin, pairs api.KVPairs) error {
	for _, kv := range pairs {

		policy, err := ParseConsulKVPair(kv)
		if err != nil {
			logging.Log.Warningf("Unable to parse key '%s' to policy: %v", kv.Key, err)
			continue
		}

		if err := UpdateNamedPolicies(plug, policy); err != nil {
			logging.Log.Warningf("Unable to update policy '%s': %v", policy.Name, err)
			continue
		}
	}

	return nil
}

func ParseConsulKVPair(kv *api.KVPair) (*policies.Policy, error) {
	reader := bytes.NewReader(kv.Value)
	policy, err := policies.ParsePolicyFile(reader)
	if err != nil {
		return nil, err
	}

	if ok, err := policy.ValidatePolicy(); !ok || err != nil {
		return nil, fmt.Errorf("unable to validate policy: %v", err)
	}

	policy.SortPolicy()
	policy.ProcessPolicyData()

	return policy, nil
}

func UpdateNamedPolicies(p *ConsulRpzPlugin, policy *policies.Policy) error {
	if policy == nil {
		return fmt.Errorf("unable to update with an empty policy")
	}

	for i := range p.Policies {
		if p.Policies[i].Name == policy.Name {

			logging.Log.Debugf("Checking hash for policy '%s':", policy.Name)
			logging.Log.Debugf("  Hash1: [%s]", p.Policies[i].Hash)
			logging.Log.Debugf("  Hash2: [%s]", policy.Hash)

			if p.Policies[i].Hash != policy.Hash {
				p.Policies[i].Disabled = policy.Disabled
				p.Policies[i].Priority = policy.Priority
				p.Policies[i].Rules = policy.Rules
				p.Policies[i].Hash = policy.Hash

				logging.Log.Debugf("Policy '%s' updated", policy.Name)
			}

			return nil
		}
	}

	logging.Log.Infof("Policy '%s' added to the list with hash ['%s']", policy.Name, policy.Hash)
	p.Policies = append(p.Policies, *policy)

	return nil
}
