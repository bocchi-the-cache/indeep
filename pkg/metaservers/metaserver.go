package metaservers

import "github.com/bocchi-the-cache/indeep/pkg/dataservers"

type PartitionID interface {
	ClientID() string
	Bucket() string
	Key() string
}

type Metadata interface {
	Size() int
	PartCount() int
	Parts() []dataservers.PartitionID

	ContentType() string
	UserMeta() map[string]string
}

type Server interface {
	Get(id PartitionID) (Partition, error)
}

type Partition interface {
	Get(id PartitionID) (Metadata, error)
	Put(id PartitionID, m Metadata) error
}
