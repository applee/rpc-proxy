package protobuf

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Constants that identify the encoding of a value on the wire.
const (
	WireVarint     = 0
	WireFixed64    = 1
	WireBytes      = 2
	WireStartGroup = 3
	WireEndGroup   = 4
	WireFixed32    = 5
)

type encoder func(*Buffer, *Field, interface{}) error
type valueEncoder func(*Buffer, uint64) error

type Type struct {
	Name      string
	KeyType   *Type // if map
	ValueType *Type // if map
}

type (
	Dependency string
	Identifier string
)

type EnumValue struct {
	Name  string
	Value int64
}

type Enum struct {
	Name   string
	Values map[string]*EnumValue
}

type Field struct {
	Name     string
	Tag      int64
	Type     *Type
	Repeated bool
	Oneof    *string
	tagcode  []byte
	enc      encoder
	valEnc   valueEncoder
}

// precalculate tag code
func (f *Field) init(root *Protobuf) error {
	var wireType uint32
	typ := f.Type.Name
	if typ == "int32" || typ == "int64" || typ == "uint32" || typ == "uint64" ||
		typ == "sint32" || typ == "sint64" || typ == "bool" || root.hasEnum(typ) {
		wireType = WireVarint
		if typ == "sint32" {
			f.valEnc = (*Buffer).EncodeZigzag32
		} else if typ == "sint64" {
			f.valEnc = (*Buffer).EncodeZigzag64
		} else {
			f.valEnc = (*Buffer).EncodeVarint
		}
	} else if typ == "fixed64" || typ == "sfixed64" || typ == "double" {
		wireType = WireFixed64
		f.valEnc = (*Buffer).EncodeFixed64
	} else if typ == "string" || typ == "bytes" || root.hasMessage(typ) {
		wireType = WireBytes
		f.valEnc = (*Buffer).EncodeVarint
	} else if typ == "fixed32" || typ == "sfixed32" || typ == "float" {
		wireType = WireFixed32
	} else if strings.HasPrefix(typ, "google.protobuf.") {
		wireType = WireBytes
		fmt.Printf("Warning: type(%s) not supported yet.\n", typ)
	} else {
		return errInvalidField
	}

	if !f.Repeated {
		switch typ {
		case "bool":
			f.enc = (*Buffer).enc_bool
		case "int32", "sint32", "sfixed32":
			f.enc = (*Buffer).enc_int32
		case "uint32", "fixed32":
			f.enc = (*Buffer).enc_uint32
		case "int64", "sint64", "sfixed64":
			f.enc = (*Buffer).enc_int64
		case "uint64", "fixed64":
			f.enc = (*Buffer).enc_int64
		case "float":
			f.enc = (*Buffer).enc_float32
		case "double":
			f.enc = (*Buffer).enc_float64
		case "string", "bytes":
			f.enc = (*Buffer).enc_string
		default:
			if strings.HasPrefix(typ, "google.protobuf.") {
				//f.enc = (*Buffer).enc_internal
			} else if root.hasEnum(typ) {
				f.enc = (*Buffer).enc_int32
			} else if root.hasMessage(typ) {
				//f.enc = (*Buffer).enc_struct_message
			} else {
				return errInvalidField
			}
		}
	} else {
		wireType = WireBytes
	}

	x := uint32(f.Tag)<<3 | wireType
	i := 0
	tagbuf := [8]byte{}
	for i = 0; x > 127; i++ {
		tagbuf[i] = 0x80 | uint8(x&0x7F)
		x >>= 7
	}
	tagbuf[i] = uint8(x)
	f.tagcode = tagbuf[0 : i+1]
	return nil
}

// validate the field
func (f *Field) validate() error {
	//TODO: check packed option
	return nil
}

type Message struct {
	Name           string
	Fields         map[string]*Field
	NestedMessages []*Message
	NestedEnums    []*Enum
	Order          []string
}

func (m *Message) init(p *Protobuf) error {
	var err error
	for _, nm := range m.NestedMessages {
		p.Messages[nm.Name] = nm
		err = nm.init(p)
		if err != nil {
			return err
		}
	}

	for _, ne := range m.NestedEnums {
		p.Enums[ne.Name] = ne
	}

	for _, f := range m.Fields {
		err = f.init(p)
		if err != nil {
			return err
		}
	}

	// order fields
	sort.Sort(m)

	return nil
}

func (m *Message) Marshal(p *Protobuf, buf *Buffer, data map[string]interface{}) error {
	for _, fn := range m.Order {
		field, ok := m.Fields[fn]
		if !ok {
			return errInvalidField
		}
		v, ok := data[field.Name]
		if !ok || v == nil {
			continue
		}
		if field.enc != nil {
			field.enc(buf, field, v)
		}
	}
	return nil
}

func (p *Message) Len() int { return len(p.Order) }
func (p *Message) Less(i, j int) bool {
	return p.Fields[p.Order[i]].Tag < p.Fields[p.Order[j]].Tag
}
func (p *Message) Swap(i, j int) { p.Order[i], p.Order[j] = p.Order[j], p.Order[i] }

type Method struct {
	Name       string
	Request    string
	Response   string
	ReqStream  bool
	RespStream bool
}

type Service struct {
	Name    string
	Methods map[string]*Method
}

type Protobuf struct {
	Syntax     string
	Dependency map[string]bool
	Enums      map[string]*Enum
	Messages   map[string]*Message
	Services   map[string]*Service
}

func (p *Protobuf) init() error {
	for _, m := range p.Messages {
		err := m.init(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Protobuf) hasEnum(name string) (ok bool) {
	_, ok = p.Enums[name]
	return
}

func (p *Protobuf) hasMessage(name string) (ok bool) {
	_, ok = p.Messages[name]
	return
}

func (p Protobuf) String() string {
	b, _ := json.MarshalIndent(p, "", " ")
	return string(b)
}

// 有效性检查
func (p *Protobuf) validate() (err error) {
	return
}
