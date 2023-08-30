package placers

import (
	"net/http"

	"github.com/bocchi-the-cache/indeep/api"
)

type client struct {
	h *http.Client
}

func (c *client) Meta(id api.MetaPartitionID) (api.MetaService, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) Data(id api.DataPartitionID) (api.DataService, error) {
	//TODO implement me
	panic("implement me")
}
