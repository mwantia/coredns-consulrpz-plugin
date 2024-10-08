package consulrpz

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
	"github.com/mwantia/coredns-consulrpz-plugin/logging"
)

type TestData struct {
	Name  string
	Match string
}

func TestRPZ(tst *testing.T) {
	OverwriteStdOut()
	clog.D.Set()

	c := caddy.NewTestController("dns", `
		consulrpz dns/tests/policies {
		  address http://127.0.0.1:8500
		  watch false
		  execution sequence
		}
	`)

	plug, err := CreatePlugin(c)
	if err != nil {
		tst.Errorf("Unable to get config: %v", err)
	}

	tests := []TestData{
		{
			Name:  "pixels.ai",
			Match: "0.0.0.0",
		},
	}

	time.Sleep(1000)
	RunTests(tst, plug, tests)
}

func RunTests(tst *testing.T, plug *ConsulRpzPlugin, tests []TestData) {
	ctx := context.TODO()

	for _, tc := range tests {
		tst.Run("Domain: "+tc.Name, func(t *testing.T) {
			logging.Log.Debugf("Testing query '%s'", tc)

			req := new(dns.Msg)
			req.SetQuestion(dns.Fqdn(tc.Name), dns.TypeA)
			rec := dnstest.NewRecorder(&test.ResponseWriter{})

			code, err := plug.ServeDNS(ctx, rec, req)

			if err != nil {
				tst.Errorf("Expected no error, but got: %v", err)
			}
			if rec.Msg == nil || len(rec.Msg.Answer) == 0 {
				tst.Errorf("Expected an answer, but got none")
			}

			answer := rec.Msg.Answer[0]
			address := answer.(*dns.A).A.String()

			logging.Log.Infof("Received code '%v', no errors", code)
			logging.Log.Infof("Answer to match '%s' with expected answer '%s'", address, tc.Match)

			if address != tc.Match {
				tst.Errorf("Expected '%s', but received '%s'", tc.Match, address)
			}
		})
	}
}

func OverwriteStdOut() error {
	tempFile, err := os.CreateTemp("", "coredns-rpz-plugin")
	if err != nil {
		return err
	}

	defer os.Remove(tempFile.Name())

	orig := logging.Log
	logging.Log = clog.NewWithPlugin("rpz")
	log.SetOutput(os.Stdout)

	defer func() {
		logging.Log = orig
	}()

	return nil
}
