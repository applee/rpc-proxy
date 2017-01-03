package protobuf

import (
	"bytes"
	"io"
	"sync"
)

type Buffer struct {
	bytes.Buffer
}

type ProtobufProtocol struct {
	bufPool sync.Pool
	parser  *Protobuf
}

func NewProtobufProtocol() (*ProtobufProtocol, error) {
	return &ProtobufProtocol{
		bufPool: sync.Pool{
			New: func() interface{} {
				return &Buffer{}
			},
		},
	}, nil
}

func (p *ProtobufProtocol) Parse(name string, reader io.Reader) error {
	pb, err := ParseReader("", reader)
	if err != nil {
		return err
	}
	p.parser = pb.(*Protobuf)
	return nil
}

func (p *ProtobufProtocol) Marshal(name string, data map[string]interface{}) ([]byte, error) {
	if p.parser == nil {
		return nil, errInvalidParser
	}
	m, ok := p.parser.Messages[name]
	if !ok || m == nil {
		return nil, &MessageNotExistsError{name}
	}

	buf := p.bufPool.Get().(*Buffer)
	defer p.bufPool.Put(buf)

	err := m.Marshal(p.parser, buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
