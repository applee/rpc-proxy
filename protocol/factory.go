package protocol

import (
	"fmt"
	"strings"

	"github.com/applee/rpc-proxy/protocol/protobuf"
	"github.com/applee/rpc-proxy/protocol/thrift"
)

type ProtocolFactory func() (Protocol, error)

var protocolFactories = make(map[string]ProtocolFactory)

func Register(name string, factory ProtocolFactory) {
	if factory == nil {
		panic(fmt.Errorf("Missing ProtocolFactory function."))
	}
	_, registered := protocolFactories[name]
	if registered {
		panic(fmt.Errorf("ProtocolFactory factory %s already registered.", name))
	}
	protocolFactories[name] = factory
}

func init() {
	Register(PROTOCOL_PB, protobuf.NewProtobufProtocol)
	Register(PROTOCOL_THRIFT, thrift.NewThriftProtocol)
}

func CreateProtocol() (MetricStore, error) {
	factory, ok := protocolFactories[name]
	if !ok {
		availableProtocols := make([]string, len(protocolFactories))
		for k, _ := range protocolFactories {
			availableProtocols = append(availableProtocols, k)
		}
		return nil, fmt.Errorf("Invalid protocol name. Must be one of: %s", strings.Join(availableProtocols, ", "))
	}

	return factory()
}
