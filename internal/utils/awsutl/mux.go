package awsutl

import (
	"net/http"
	"strings"
)

type S3Router struct {
	ListBuckets http.HandlerFunc
}

func (r *S3Router) Route(req *http.Request) http.HandlerFunc {
	var bucketName, objectName string
	{
		parts := strings.SplitN(strings.TrimLeft(req.URL.Path, "/"), "/", 2)
		switch len(parts) {
		case 1:
			bucketName = parts[0]
		case 2:
			bucketName = parts[0]
			objectName = parts[1]
		}
	}
	if objectName == "" && bucketName == "" {
		return r.routeMyOperations(req)
	}
	return nil
}

func (r *S3Router) routeMyOperations(req *http.Request) http.HandlerFunc {
	switch req.Method {
	case http.MethodGet:
		return r.ListBuckets
	default:
		return nil
	}
}
