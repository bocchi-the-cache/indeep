package api

const DefaultGatewayHost = "127.0.0.1:11451"

type ListAllMyBucketsResult struct{}

type Gateway interface {
	ListBuckets() (*ListAllMyBucketsResult, error)
}
