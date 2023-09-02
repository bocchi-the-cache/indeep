package typedh

import "net/http"

type (
	ConsumerFunc[Req any]        func(r *Req) error
	ProviderFunc[Resp any]       func() (*Resp, error)
	ProcessorFunc[Req, Resp any] func(r *Req) (*Resp, error)
)

func ConsumerWith[Req any](c Codec, f ConsumerFunc[Req]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(c, w, r)
		var req Req
		if err := UnmarshalWith(c, r.Body, &req); err != nil {
			ctx.Err(err)
			return
		}
		if err := f(&req); err != nil {
			ctx.Err(err)
			return
		}
	}
}

func ProviderWith[Resp any](c Codec, f ProviderFunc[Resp]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(c, w, r)
		resp, err := f()
		if err != nil {
			ctx.Err(err)
			return
		}
		ctx.OK(resp)
	}
}

func ProcessorWith[Req, Resp any](c Codec, f ProcessorFunc[Req, Resp]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Req
		ctx := NewContext(c, w, r)
		if err := UnmarshalWith(c, r.Body, &req); err != nil {
			ctx.Err(err)
			return
		}
		resp, err := f(&req)
		if err != nil {
			ctx.Err(err)
			return
		}
		ctx.OK(resp)
	}
}

func Consumer[Req any](f ConsumerFunc[Req]) http.HandlerFunc   { return ConsumerWith(DefaultCodec, f) }
func Provider[Resp any](f ProviderFunc[Resp]) http.HandlerFunc { return ProviderWith(DefaultCodec, f) }
func Processor[Req, Resp any](f ProcessorFunc[Req, Resp]) http.HandlerFunc {
	return ProcessorWith(DefaultCodec, f)
}
