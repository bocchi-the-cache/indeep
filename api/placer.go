package api

type Placer interface {
	Member

	LookupMetaService(key MetaKey) (MetaService, error)
	AddMetaService( /* TODO */ ) error

	LookupDataService(id DataPartitionID) (DataService, error)
	AddDataService( /* TODO */ ) error
}
