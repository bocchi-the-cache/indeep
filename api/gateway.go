package api

import "net/http"

const DefaultGatewayHost = "127.0.0.1:11451"

type ListAllMyBucketsResult struct{}

type Gateway interface {
	ListBuckets() (*ListAllMyBucketsResult, error)
}

type S3ApiID int

const (
	ListBuckets S3ApiID = iota

	S3ApiMaxSize
)

type S3Mux interface {
	http.Handler
	HandleFunc(id S3ApiID, f http.HandlerFunc)
}
