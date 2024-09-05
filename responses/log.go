package responses

import (
	"encoding/json"
	"os"
	"time"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/mwantia/coredns-consulrpz-plugin/matches"
	"github.com/mwantia/coredns-consulrpz-plugin/policies"
)

type LogMessageEntry struct {
	Time     string `json:"time"`
	QName    string `json:"qname"`
	QType    uint16 `json:"qtype"`
	RemoteIP string `json:"remoteip"`
	Result   string `json:"result"`
}

func HandleLogResponse(state request.Request, value json.RawMessage, result *matches.MatchResult, policy policies.Policy, response *PolicyResponse) error {
	var path string
	if err := json.Unmarshal(value, &path); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	entry := LogMessageEntry{
		Time:     time.Now().Format("2006-01-02 15:04:05"),
		QName:    dns.Fqdn(state.Name()),
		QType:    state.QType(),
		RemoteIP: state.IP(),
		Result:   result.Data,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, err = f.WriteString(string(data) + "\n")
	return err
}
