package api

type Placer interface {
	Member

	LookupMetaClient(key MetaKey) (MetaService, error)
	AddMetaServer( /* TODO */ ) error

	LookupDataClient(id DataPartitionID) (DataService, error)
	AddDataServer( /* TODO */ ) error
}
