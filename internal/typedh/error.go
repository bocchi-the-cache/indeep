package typedh

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
	ErrorData interface{ ErrorData() any }

	withData struct {
		error
		v any
	}
)

func WithData(err error, v any) error { return &withData{error: err, v: v} }
func (w *withData) ErrorData() any    { return w.v }
