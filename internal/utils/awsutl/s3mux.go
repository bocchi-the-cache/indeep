package awsutl

import (
	"net/http"
	"strings"

	"github.com/bocchi-the-cache/indeep/api"
)

type S3Mux struct {
	handlers [api.S3ApiMaxSize]http.HandlerFunc
}

func (s *S3Mux) HandleFunc(id api.S3ApiID, f http.HandlerFunc) { s.handlers[id] = f }

func (s *S3Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if id, ok := s.parseApiID(r); ok {
		s.handlers[id].ServeHTTP(w, r)
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *S3Mux) parseApiID(r *http.Request) (api.S3ApiID, bool) {
	const maxSeps = 3
	var objectName, bucketName string
	{
		parts := strings.SplitN(r.URL.Path, "/", maxSeps)
		switch len(parts) {
		case maxSeps - 1:
			objectName = parts[1]
		case maxSeps:
			objectName = parts[1]
			bucketName = parts[2]
		default:
			return 0, false
		}
	}
	if objectName == "" && bucketName == "" {
		return api.ListBuckets, true
	}
	return 0, false
}
