package api

const (
	RpcGetMembers        = "get-members"
	RpcAskLeader         = "ask-leader"
	RpcLookupMetaService = "lookup-meta-service"
	RpcAddMetaService    = "add-meta-service"
	RpcLookupDataService = "lookup-data-service"
	RpcAddDataService    = "add-data-service"
)

type Placer interface {
	Member

	LookupMetaService(key MetaKey) (MetaService, error)
	AddMetaService( /* TODO */ ) error

	LookupDataService(id DataPartitionID) (DataService, error)
	AddDataService( /* TODO */ ) error
}
