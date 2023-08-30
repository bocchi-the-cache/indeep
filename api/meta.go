package api

type MetaKey interface {
	ClientID() string
	Bucket() string
	Key() string
}

type MetaPartition interface {
	StartKey() string
	EndKey() string
	KeyCount() int

	Get(key string) (Metadata, error)
	Put(key string, m Metadata) error
}

type Metadata interface {
	Size() int
	PartCount() int
	Parts() []DataPartitionID

	ContentType() string
	UserMeta() map[string]string
}

type MetaService interface {
	Lookup(key MetaKey) (MetaPartition, error)
}
