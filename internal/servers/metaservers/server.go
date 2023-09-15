package metaservers

import (
	"context"

	"github.com/bocchi-the-cache/indeep/api"
)

func (m *metaserver) ListenAndServe() error              { return m.server.ListenAndServe() }
func (m *metaserver) Shutdown(ctx context.Context) error { return m.server.Shutdown(ctx) }

func (m *metaserver) Lookup(key api.MetaKey) (api.MetaPartition, error) {
	//TODO implement me
	panic("implement me")
}
