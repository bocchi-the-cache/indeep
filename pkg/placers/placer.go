package placers

import (
	"github.com/bocchi-the-cache/indeep/pkg/dataservers"
	"github.com/bocchi-the-cache/indeep/pkg/metaservers"
)

type Placer interface {
	Meta(id metaservers.PartitionID) (metaservers.Server, error)
	Data(id dataservers.PartitionID) (dataservers.Server, error)
}
