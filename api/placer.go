package api

import "github.com/hashicorp/raft"

const (
	DefaultPlacerID                      = "placer0"
	DefaultPlacerHost                    = "127.0.0.1:11551"
	DefaultPlacerPeer raft.ServerAddress = "127.0.0.1:11561"

	RpcGetMembers        RpcID = "get-members"
	RpcAskLeader         RpcID = "ask-leader"
	RpcLookupMetaService RpcID = "lookup-meta-service"
	RpcAddMetaService    RpcID = "add-meta-service"
	RpcLookupDataService RpcID = "lookup-data-service"
	RpcAddDataService    RpcID = "add-data-service"
)

type Placer interface {
	Member

	LookupMetaService(key MetaKey) (MetaService, error)
	AddMetaService( /* TODO */ ) error

	LookupDataService(id DataPartitionID) (DataService, error)
	AddDataService( /* TODO */ ) error
}
