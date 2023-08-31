package jsonhttp

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/bocchi-the-cache/indeep/internal/logs"
)

func Unmarshal(body io.ReadCloser, v any) error {
	defer func() { _ = body.Close() }()
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

type (
	StatusCoder interface{ StatusCode() int }

	withStatusCode struct {
		error
		c int
	}
)

func WithStatusCode(err error, c int) error { return &withStatusCode{error: err, c: c} }
func (w *withStatusCode) StatusCode() int   { return w.c }

type (
	JSONResponseWriter interface {
		OK(v any)
		Err(err error)
	}

	jsonRespWriter struct{ w http.ResponseWriter }
)

func W(w http.ResponseWriter) JSONResponseWriter { return &jsonRespWriter{w} }

func (j *jsonRespWriter) OK(v any) {
	data, err := json.Marshal(v)
	if err != nil {
		j.Err(err)
		return
	}
	j.write(http.StatusOK, data)
}

func (j *jsonRespWriter) Err(err error) {
	statusCode := http.StatusInternalServerError
	if c, ok := err.(StatusCoder); ok {
		statusCode = c.StatusCode()
	}
	data, _ := json.Marshal(struct{ Msg string }{err.Error()})
	j.write(statusCode, data)
}

func (j *jsonRespWriter) write(statusCode int, data []byte) {
	j.w.Header().Add("Content-Type", "application/json")
	j.w.WriteHeader(statusCode)
	if _, err := j.w.Write(data); err != nil {
		logs.E.Println("write error:", err)
	}
}
