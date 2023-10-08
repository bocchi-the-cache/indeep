package hyped

import (
	"io"
	"net/http"

	"github.com/bocchi-the-cache/indeep/internal/logs"
)

type (
	Context interface {
		OK(v any)
		Err(err error)
	}

	context struct {
		c Codec
		w http.ResponseWriter
		r *http.Request
	}
)

func NewContext(c Codec, w http.ResponseWriter, r *http.Request) Context {
	return &context{c: c, w: w, r: r}
}

func unmarshalBody(c Codec, body io.ReadCloser, v any) error {
	defer func() { _ = body.Close() }()
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	return c.Unmarshal(data, v)
}

func (c *context) OK(v any) {
	data, err := c.c.Marshal(v)
	if err != nil {
		c.Err(err)
		return
	}
	c.write(http.StatusOK, data)
}

func (c *context) Err(err error) {
	statusCode := http.StatusInternalServerError
	if c, ok := err.(StatusCoder); ok {
		statusCode = c.StatusCode()
	}
	data, marshalErr := c.c.Marshal(err)
	if marshalErr != nil {
		c.write(http.StatusInternalServerError, nil)
		logs.S.Error("failed to marshal error response", "err", marshalErr)
		return
	}
	c.write(statusCode, data)
}

func (c *context) write(statusCode int, data []byte) {
	c.w.Header().Add("Content-Type", c.c.ContentType())
	c.w.WriteHeader(statusCode)
	if data == nil {
		return
	}
	if _, err := c.w.Write(data); err != nil {
		logs.S.Error("write error", "err", err)
	}
}
