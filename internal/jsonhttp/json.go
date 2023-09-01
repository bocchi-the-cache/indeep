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
	ResponseWriter interface {
		OK(v any)
		Err(err error)
	}

	respWriter struct{ w http.ResponseWriter }
)

func W(w http.ResponseWriter) ResponseWriter { return &respWriter{w} }

func (w *respWriter) OK(v any) {
	data, err := json.Marshal(v)
	if err != nil {
		w.Err(err)
		return
	}
	w.write(http.StatusOK, data)
}

func (w *respWriter) Err(err error) {
	statusCode := http.StatusInternalServerError
	if c, ok := err.(StatusCoder); ok {
		statusCode = c.StatusCode()
	}
	data, _ := json.Marshal(struct{ Msg string }{err.Error()})
	w.write(statusCode, data)
}

func (w *respWriter) write(statusCode int, data []byte) {
	w.w.Header().Add("Content-Type", "application/json")
	w.w.WriteHeader(statusCode)
	if _, err := w.w.Write(data); err != nil {
		logs.E.Println("write error:", err)
	}
}
