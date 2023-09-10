package groups

import (
	"net/url"

	"github.com/bocchi-the-cache/indeep/api"
)

type multipeer struct{}

func (m *multipeer) String() string {
	//TODO implement me
	panic("implement me")
}

func (m *multipeer) URL() *url.URL {
	//TODO implement me
	panic("implement me")
}

func (m *multipeer) RPC(id api.RpcID) *url.URL {
	//TODO implement me
	panic("implement me")
}

func (m *multipeer) MarshalJSON() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (m *multipeer) UnmarshalJSON(bytes []byte) error {
	//TODO implement me
	panic("implement me")
}
