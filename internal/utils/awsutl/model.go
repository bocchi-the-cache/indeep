package awsutl

type Error struct {
	Status  int    `xml:"-"`
	Code    string `xml:",omitempty"`
	Message string `xml:",omitempty"`
}

func (e *Error) Error() string   { return e.Message }
func (e *Error) StatusCode() int { return e.Status }
