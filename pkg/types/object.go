package types

import (
	"github.com/bocchi-the-cache/indeep/pkg/metaservers"
)

type Object interface {
	metaservers.PartitionID
	metaservers.Metadata
}
