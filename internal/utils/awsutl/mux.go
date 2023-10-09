package awsutl

import (
	"net/http"
	"strings"
)

type S3Mux struct {
	ListBuckets http.HandlerFunc
}

func (m *S3Mux) Route(r *http.Request) http.HandlerFunc {
	var bucketName, objectName string
	{
		parts := strings.SplitN(strings.TrimLeft(r.URL.Path, "/"), "/", 2)
		switch len(parts) {
		case 1:
			bucketName = parts[0]
		case 2:
			bucketName = parts[0]
			objectName = parts[1]
		}
	}
	if objectName == "" && bucketName == "" {
		return m.routeMyOperations(r)
	}
	return nil
}

func (m *S3Mux) routeMyOperations(r *http.Request) http.HandlerFunc {
	switch r.Method {
	case http.MethodGet:
		return m.ListBuckets
	default:
		return nil
	}
}
