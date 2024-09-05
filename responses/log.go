package responses

import (
	"bytes"
	"encoding/json"
	"os"
	"text/template"
	"time"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/mwantia/coredns-consulrpz-plugin/matches"
	"github.com/mwantia/coredns-consulrpz-plugin/policies"
)

var DefaultLogFormat = "{{.Time}} [consulrpz] {{.RemoteIP}} made {{.QType}} query for {{.QName}}"

type LogResponse struct {
	Path   string `json:"path"`
	AsJson bool   `json:"as_json"`
	Format string `json:"format"`
}

type LogMessage struct {
	Time     string `json:"time"`
	QName    string `json:"qname"`
	QType    uint16 `json:"qtype"`
	RemoteIP string `json:"remoteip"`
	Result   string `json:"result"`
}

func HandleLogResponse(state request.Request, value json.RawMessage, result *matches.MatchResult, policy policies.Policy, response *PolicyResponse) error {
	var log LogResponse
	if err := json.Unmarshal(value, &log); err != nil {
		return err
	}

	f, err := os.OpenFile(log.Path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	msg := LogMessage{
		Time:     time.Now().Format("2006-01-02 15:04:05"),
		QName:    dns.Fqdn(state.Name()),
		QType:    state.QType(),
		RemoteIP: state.IP(),
		Result:   result.Data,
	}

	if len(log.Format) <= 0 {
		log.Format = DefaultLogFormat
	}

	if log.AsJson {
		return LogJsonResponse(*f, msg)
	}

	return LogFormatResponse(*f, msg, log.Format)
}

func LogFormatResponse(f os.File, msg LogMessage, format string) error {
	tmpl, err := template.New(f.Name()).Parse(format)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, msg)
	if err != nil {
		return err
	}

	text := buffer.String() + "\n"
	_, err = f.WriteString(text)

	return err
}

func LogJsonResponse(f os.File, msg LogMessage) error {
	json, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = f.WriteString(string(json) + "\n")

	return err
}
