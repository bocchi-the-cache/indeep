package awsutl

import (
	"net/http"
	"strings"

	"github.com/bocchi-the-cache/indeep/api"
)

type s3mux struct {
	handlers [api.S3ApiMaxSize]http.HandlerFunc
}

func S3Mux() api.S3Mux { return new(s3mux) }

func (s *s3mux) HandleFunc(id api.S3ApiID, f http.HandlerFunc) { s.handlers[id] = f }

func (s *s3mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if id, ok := s.parseApiID(r); ok {
		s.handlers[id].ServeHTTP(w, r)
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *s3mux) parseApiID(r *http.Request) (api.S3ApiID, bool) {
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
