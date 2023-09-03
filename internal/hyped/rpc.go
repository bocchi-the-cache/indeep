package hyped

import (
	"net/http"

	"github.com/bocchi-the-cache/indeep/api"
)

type (
	RPC interface {
		Get(p api.Peer, id api.RpcID, v any) error
	}

	rpc struct {
		h *http.Client
		c Codec
	}
)

func NewRPCWithCodec(h *http.Client, c Codec) RPC { return &rpc{h: h, c: c} }
func NewRPC(h *http.Client) RPC                   { return NewRPCWithCodec(h, DefaultCodec) }

func (r *rpc) Get(p api.Peer, id api.RpcID, v any) error {
	resp, err := r.h.Get(p.RPC(id).String())
	if err != nil {
		return err
	}
	return unmarshalBody(r.c, resp.Body, v)
}
