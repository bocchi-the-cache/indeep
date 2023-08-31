package api

import (
	"fmt"
	"net/url"
)

type Endpoint interface {
	fmt.Stringer
	URL() *url.URL
}
