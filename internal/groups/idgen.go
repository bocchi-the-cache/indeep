package groups

import (
	"fmt"
	"sync/atomic"

	"github.com/bocchi-the-cache/indeep/api"
)

const GroupIDPrefix = "g"

type IDGen interface {
	New() api.GroupID
}

type idGen struct {
	c atomic.Uint64
}

func NewIDGen(init uint64) IDGen {
	g := new(idGen)
	g.c.Store(init)
	return g
}

func (i *idGen) New() api.GroupID {
	id := i.c.Add(1)
	return api.GroupID(fmt.Sprint(GroupIDPrefix, id))
}
