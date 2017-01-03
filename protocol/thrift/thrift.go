package thrift

type ThriftProtocol struct {
	Protocol string
}

func NewThriftProtocol() (*ThriftProtocol, error) {
	return &ThriftProtocol{}, nil
}
