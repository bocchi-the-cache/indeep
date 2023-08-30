package api

import "encoding"

type MetaPartitionID interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	ClientID() string
	Bucket() string
	Key() string
}

type Metadata interface {
	Size() int
	PartCount() int
	Parts() []DataPartitionID

	ContentType() string
	UserMeta() map[string]string
}

type MetaService interface {
	Get(id MetaPartitionID) (MetaPartition, error)
}

type MetaPartition interface {
	Get(id MetaPartitionID) (Metadata, error)
	Put(id MetaPartitionID, m Metadata) error
}
