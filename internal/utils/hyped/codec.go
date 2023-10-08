package hyped

import (
	"encoding/json"
	"encoding/xml"
)

type Codec interface {
	ContentType() string
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

var DefaultCodec = JSON()

type jsonEncoding struct{}

func JSON() Codec                                        { return new(jsonEncoding) }
func (*jsonEncoding) ContentType() string                { return "application/json; charset=utf-8" }
func (*jsonEncoding) Marshal(v any) ([]byte, error)      { return json.Marshal(v) }
func (*jsonEncoding) Unmarshal(data []byte, v any) error { return json.Unmarshal(data, v) }

type xmlEncoding struct{}

func XML() Codec                                        { return new(xmlEncoding) }
func (*xmlEncoding) ContentType() string                { return "application/xml; charset=utf-8" }
func (*xmlEncoding) Marshal(v any) ([]byte, error)      { return xml.Marshal(v) }
func (*xmlEncoding) Unmarshal(data []byte, v any) error { return xml.Unmarshal(data, v) }
