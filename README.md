# CoreDNS RPZ Plugin

This plugin enables CoreDNS to use custom Response Policy Zones (RPZ) for DNS filtering and policy enforcement.

> [!IMPORTANT]
> This plugin is still actively being worked on. \
> Expect possible changes or reworks of how this plugin functions and how the config is structured.
>
> Additionally, this README isn't always up-to-date, so not everything mentioned here might work as described.

## Features

- Use Consul KV as a backend for RPZ policies
- Real-time policy updates via Consul KV
- Support for various RPZ triggers and actions
- Configurable policy priorities
- Metrics for monitoring (compatible with Prometheus)

## Architecture

The CoreDNS RPZ Plugin follows a modular architecture to process DNS queries and apply RPZ policies:

```mermaid
sequenceDiagram
    participant Client
    participant CoreDNS
    participant RPZPlugin
    participant ConsulKV
    participant PolicyHandler
    participant TriggerHandler
    participant ActionHandler

    Client->>CoreDNS: DNS Query
    CoreDNS->>RPZPlugin: ServeDNS()
    RPZPlugin->>ConsulKV: Fetch Policies
    ConsulKV-->>RPZPlugin: Return Policies
    RPZPlugin->>PolicyHandler: HandlePoliciesParallel()
    loop For each Policy
        PolicyHandler->>TriggerHandler: Check Triggers
        alt Triggers Match
            TriggerHandler-->>PolicyHandler: Triggers Matched
            PolicyHandler->>ActionHandler: Execute Actions
            ActionHandler-->>PolicyHandler: Action Result
        else Triggers Don't Match
            TriggerHandler-->>PolicyHandler: No Match
        end
    end
    PolicyHandler-->>RPZPlugin: Policy Result
    alt Policy Applied
        RPZPlugin-->>CoreDNS: Modified DNS Response
    else No Policy Match
        RPZPlugin-->>CoreDNS: Original DNS Query
    end
    CoreDNS-->>Client: DNS Response
```

## Installation

To use this plugin, you need to compile it into CoreDNS. Add the following line to the `plugin.cfg` file in your CoreDNS source code:

```
rpz:github.com/mwantia/coredns-rpz-plugin
```

Then, rebuild CoreDNS with:

```sh
go get github.com/mwantia/coredns-rpz-plugin
go generate
go build
```

## Configuration

Add the plugin to your CoreDNS configuration file (Corefile):

```corefile
rpz {
    file policies/example.json
    consul dns/policies
}
```

### Configuration options:

- `policy`: Specifies a json-file used for storing RPZ policies 
- `consul`: Specifies the Consul KV prefix for storing RPZ policies

## RPZ Configuration

These custom RPZ policies are written in JSON. Each policy should be a object with the following structure:

```json
{
    "name": "RPZ Example",
    "version": "1.0",
    "priority": 0,
    "rules": [
        {
            "triggers": [
                {
                    "type": "domain",
                    "value": ["example.com"]
                }
            ],
            "actions": [
                {
                    "type": "deny"
                }
            ]
        }
    ]
}
```

### Policy structure:
- `name`: Name of the policy (string)
- `version`: Version of the policy format (string, must be "1.0")
- `priority`: Priority of the policy (integer, lower values have higher priority; Default `1000`)
- `rules`: Array of policy rules

### Rule structure:
- `triggers`: Array of conditions that trigger the rule
- `actions`: Array of actions to take when the rule is triggered

## Supported Triggers

Currently, the plugin supports the following trigger types:

1. `type`: Matches query types
  ```json
  {
    "type": "type",
    "value": ["A", "AAAA"]
  }
  ```

2. `name`: Matches domain names as suffix
  ```json
  {
    "type": "name",
    "value": ["example.com", "www.site.com"]
  }
  ```

3. `cidr`: Matches IP-Adress ranges
  ```json
  {
    "type": "cidr",
    "value": ["192.168.0.0/16", "192.168.178.1/32", "192.168.0.1"]
  }
  ```

4. `time`: Matches time-frames with a start and end
  ```json
  {
    "type": "time",
    "value": [
      {
        "start": "09:00",
        "end": "17:00"
      }
    ]
  }
  ```

5. `cron`: Matches time-frames declared as cron
  ```json
  {
    "type": "cron",
    "value": ["* 9-16 * * 1-5"]
  }
  ```

5. `regex`: Matches domain names via regex
  ```json
  {
    "type": "regex",
    "value": ["(.*)example\\.com"]
  }
  ```

Triggers are handled in the following order: `type`, `name`, `cidr`, `time`, `cron`, `regex`.

## Supported Actions

The plugin supports the following action types:

1. `deny`: Denies the request
   ```json
   {
     "type": "deny"
   }
   ```

2. `fallthrough`: Continues to the next plugin
   ```json
   {
     "type": "fallthrough"
   }
   ```

3. `code`: Returns a specific DNS response code
   ```json
   {
     "type": "code",
     "value": "NXDOMAIN"
   }
   ```

4. `record`: Returns a specific DNS record
   ```json
   {
     "type": "record",
     "value": {
       "ttl": 3600,
       "records": [
         {
           "type": "A",
           "value": ["0.0.0.0"]
         }
       ]
     }
   }
   ```

Actions are handled in the following order: `deny`, `fallthrough`, `code`, `record`.

## Additional TXT

```
$ dig example.com

; <<>> DiG 9.18.28-0ubuntu0.22.04.1-Ubuntu <<>> example.com
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: REFUSED, id: 59853
;; flags: qr aa rd; QUERY: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 2
;; WARNING: recursion requested but not available

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 1232
; COOKIE: a3c6576853e116eb (echoed)
;; QUESTION SECTION:
;example.com.       IN      A

;; ADDITIONAL SECTION:
example.com. 300    IN      TXT     "Handled by RPZ policy - RPZ Example"

;; Query time: 40 msec
;; SERVER: 127.0.0.1#53(127.0.0.1) (UDP)
;; WHEN: Fri Aug 30 22:40:51 CEST 2024
;; MSG SIZE  rcvd: 149
```

## Metrics

This plugin exposes the following metrics for Prometheus:

* `coredns_rpz_request_duration_seconds{status, le}`: 
  * Histogram of the time (in seconds) each request to Consul took
* `coredns_rpz_query_requests_total{status, policy, type}`:
  * Count of the queries received and processed by the plugin

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
