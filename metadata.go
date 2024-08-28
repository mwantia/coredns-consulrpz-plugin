package rpz

import (
	"context"
	"strings"

	"github.com/coredns/coredns/plugin/metadata"
	"github.com/coredns/coredns/request"
)

const MetadataRpzQueryStatus = "rpz/query-status"
const DefaultQueryStatus = QueryStatusNoMatch

func (p RpzPlugin) Metadata(ctx context.Context, state request.Request) context.Context {
	metadata.SetValueFunc(ctx, MetadataRpzQueryStatus, func() string {
		return DefaultQueryStatus
	})
	return ctx
}

func (p RpzPlugin) SetMetadataQueryStatus(ctx context.Context, status string) {
	s := strings.ToUpper(status)
	metadata.SetValueFunc(ctx, MetadataRpzQueryStatus, func() string {
		return s
	})
}
