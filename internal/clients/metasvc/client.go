package metasvc

import (
	"github.com/bocchi-the-cache/indeep/api"
)

type client struct{}

func (c *client) Get(id api.MetaPartitionID) (api.MetaPartition, error) {
	//TODO implement me
	panic("implement me")
}
