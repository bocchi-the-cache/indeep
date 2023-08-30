package api

type Placer interface {
	Meta(id MetaPartitionID) (MetaService, error)
	Data(id DataPartitionID) (DataService, error)
}
