package api

const (
	DefaultMetaserverHost      = "127.0.0.1:11651"
	DefaultMetaServerMultiPeer = "127.0.0.1:11661"
)

type MetaKey interface {
	ClientID() string
	Bucket() string
	Key() string
}

type MetaPartition interface {
	Peer

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
