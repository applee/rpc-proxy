package protocol

import "io"

const (
	PROTOCOL_PB     = "protobuf"
	PROTOCOL_THRIFT = "thrift"
)

type Protocol interface {
	Parse(string, io.Reader)
	Marshal()
	Unmarshal()
}
