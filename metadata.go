package consulrpz

import (
	"context"
	"strings"

	"github.com/coredns/coredns/plugin/metadata"
	"github.com/coredns/coredns/request"
	"github.com/mwantia/coredns-consulrpz-plugin/metrics"
)

const MetadataRpzQueryStatus = "consulrpz/query-status"
const DefaultQueryStatus = metrics.QueryStatusNoMatch

func (p ConsulRpzPlugin) Metadata(ctx context.Context, state request.Request) context.Context {
	metadata.SetValueFunc(ctx, MetadataRpzQueryStatus, func() string {
		return DefaultQueryStatus
	})
	return ctx
}

func (p ConsulRpzPlugin) SetMetadataQueryStatus(ctx context.Context, status string) {
	s := strings.ToUpper(status)
	metadata.SetValueFunc(ctx, MetadataRpzQueryStatus, func() string {
		return s
	})
}
