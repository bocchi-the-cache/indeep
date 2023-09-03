package hyped

import "encoding/json"

type Codec interface {
	ContentType() string
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

var DefaultCodec = JSON()

type jsonEncoding struct{}

func JSON() Codec                                        { return new(jsonEncoding) }
func (*jsonEncoding) ContentType() string                { return "application/json" }
func (*jsonEncoding) Marshal(v any) ([]byte, error)      { return json.Marshal(v) }
func (*jsonEncoding) Unmarshal(data []byte, v any) error { return json.Unmarshal(data, v) }
