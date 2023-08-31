package jsonutl

import (
	"encoding/json"
	"io"
)

func UnmarshalBody(r io.ReadCloser, v any) error {
	defer func() { _ = r.Close() }()
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
