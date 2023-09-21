package api

const (
	DefaultPlacerID   = "placer0"
	DefaultPlacerHost = "127.0.0.1:11551"
	DefaultPlacerPeer = "127.0.0.1:11561"

	RpcPlacerListGroups = "list-groups"
)

var (
	DefaultPlacerHostMap = NewAddressMap(HostScheme).Join(DefaultPlacerID, DefaultPlacerHost)
	DefaultPlacerPeerMap = NewAddressMap(RaftScheme).Join(DefaultPlacerID, DefaultPlacerPeer)
)

type Placer interface {
	Member

	ListGroups() ([]GroupID, error)
}
