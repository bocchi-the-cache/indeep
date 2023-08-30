package api

import "encoding"

type DataPartitionID interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

type DataService interface{}

type DataPartition interface{}
