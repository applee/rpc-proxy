// partially copied from https://github.com/golang/protobufprotobuf

package protobuf

/*
 * Routines for encoding data into the wire format for protocol buffers.
 */

import (
	"math"

	"github.com/spf13/cast"
)

// The fundamental encoders that put bytes on the wire.
// Those that take integer types all accept uint64 and are
// therefore of type valueEncoder.

const maxVarintBytes = 10 // maximum length of a varint

// maxMarshalSize is the largest allowed size of an encoded protobuf,
// since C++ and Java use signed int32s for the size.
const maxMarshalSize = 1<<31 - 1

const (
	maxInt32 = 1<<31 - 1
	minInt32 = -maxInt32 - 1
)

// EncodeVarint writes a varint-encoded integer to the Buffer.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func (p *Buffer) EncodeVarint(x uint64) error {
	for x >= 1<<7 {
		p.WriteByte(uint8(x&0x7f | 0x80))
		x >>= 7
	}
	p.WriteByte(uint8(x))
	return nil
}

// EncodeFixed64 writes a 64-bit integer to the Buffer.
// This is the format for the
// fixed64, sfixed64, and double protocol buffer types.
func (p *Buffer) EncodeFixed64(x uint64) error {
	p.Write([]byte{
		uint8(x),
		uint8(x >> 8),
		uint8(x >> 16),
		uint8(x >> 24),
		uint8(x >> 32),
		uint8(x >> 40),
		uint8(x >> 48),
		uint8(x >> 56),
	})
	return nil
}

// EncodeFixed32 writes a 32-bit integer to the Buffer.
// This is the format for the
// fixed32, sfixed32, and float protocol buffer types.
func (p *Buffer) EncodeFixed32(x uint64) error {
	p.Write([]byte{
		uint8(x),
		uint8(x >> 8),
		uint8(x >> 16),
		uint8(x >> 24),
	})
	return nil
}

// EncodeZigzag64 writes a zigzag-encoded 64-bit integer
// to the Buffer.
// This is the format used for the sint64 protocol buffer type.
func (p *Buffer) EncodeZigzag64(x uint64) error {
	// use signed number to get arithmetic right shift.
	return p.EncodeVarint(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}

// EncodeZigzag32 writes a zigzag-encoded 32-bit integer
// to the Buffer.
// This is the format used for the sint32 protocol buffer type.
func (p *Buffer) EncodeZigzag32(x uint64) error {
	// use signed number to get arithmetic right shift.
	return p.EncodeVarint(uint64((uint32(x) << 1) ^ uint32((int32(x) >> 31))))
}

// EncodeRawBytes writes a count-delimited byte buffer to the Buffer.
// This is the format used for the bytes protocol buffer
// type and for embedded messages.
func (p *Buffer) EncodeRawBytes(b []byte) error {
	p.EncodeVarint(uint64(len(b)))
	p.Write(b)
	return nil
}

// EncodeStringBytes writes an encoded string to the Buffer.
// This is the format used for the proto2 string type.
func (p *Buffer) EncodeStringBytes(s string) error {
	p.EncodeVarint(uint64(len(s)))
	p.WriteString(s)
	return nil
}

// Encode a bool.
func (o *Buffer) enc_bool(f *Field, i interface{}) error {
	v, err := cast.ToBoolE(i)
	if err != nil {
		return errInvalidDataType
	}
	if !v {
		return ErrNil
	}
	o.Write(f.tagcode)
	f.valEnc(o, 1)
	return nil
}

// Encode an int32.
func (o *Buffer) enc_int32(f *Field, i interface{}) error {
	x, err := cast.ToIntE(i)
	if err != nil {
		return errInvalidDataType
	}
	if x == 0 {
		return ErrNil
	}
	o.Write(f.tagcode)
	f.valEnc(o, uint64(x))
	return nil
}

func (o *Buffer) enc_float32(f *Field, i interface{}) error {
	v, err := cast.ToFloat64E(i)
	if err != nil {
		return errInvalidDataType
	}
	x := math.Float32bits(float32(v))
	if x == 0 {
		return ErrNil
	}
	o.Write(f.tagcode)
	f.valEnc(o, uint64(x))
	return nil
}

// Encode an double.
func (o *Buffer) enc_float64(f *Field, i interface{}) error {
	v, err := cast.ToFloat64E(i)
	if err != nil {
		return errInvalidDataType
	}
	x := math.Float64bits(v)
	if x == 0 {
		return ErrNil
	}
	o.Write(f.tagcode)
	f.valEnc(o, x)
	return nil
}

// Encode an uint32.
func (o *Buffer) enc_uint32(f *Field, i interface{}) error {
	x, err := cast.ToIntE(i)
	if err != nil {
		return errInvalidDataType
	}
	if x == 0 {
		return ErrNil
	}
	o.Write(f.tagcode)
	f.valEnc(o, uint64(x))
	return nil
}

// Encode an int64.
func (o *Buffer) enc_int64(f *Field, i interface{}) error {
	x, err := cast.ToInt64E(i)
	if err != nil {
		return errInvalidDataType
	}
	if x == 0 {
		return ErrNil
	}
	o.Write(f.tagcode)
	f.valEnc(o, uint64(x))
	return nil
}

// Encode a string.
func (o *Buffer) enc_string(f *Field, i interface{}) error {
	v, err := cast.ToStringE(i)
	if err != nil {
		return errInvalidDataType
	}
	if v == "" {
		return ErrNil
	}
	o.Write(f.tagcode)
	o.EncodeStringBytes(v)
	return nil
}

// Encode a slice of bools ([]bool) in packed format.
func (o *Buffer) enc_slice_bool(f *Field, i interface{}) error {
	s, err := cast.ToBoolSliceE(i)
	if err != nil {
		return errInvalidDataType
	}
	l := len(s)
	if l == 0 {
		return ErrNil
	}
	o.Write(f.tagcode)
	o.EncodeVarint(uint64(l)) // each bool takes exactly one byte
	for _, x := range s {
		v := uint64(0)
		if x {
			v = 1
		}
		f.valEnc(o, v)
	}
	return nil
}

// Encode a slice of bytes ([]byte).
func (o *Buffer) enc_slice_byte(f *Field, i interface{}) error {
	s, err := ToByteSliceE(i)
	if err != nil {
		return errInvalidDataType
	}
	if len(s) == 0 {
		return ErrNil
	}
	o.Write(f.tagcode)
	o.EncodeRawBytes(s)
	return nil
}

// Encode a slice of int32s ([]int32) in packed format.
func (o *Buffer) enc_slice_int32(f *Field, i interface{}) error {
	s, err := cast.ToIntSliceE(i)
	if err != nil {
		return errInvalidDataType
	}
	l := len(s)
	if l == 0 {
		return ErrNil
	}

	buf := &Buffer{}
	for i := 0; i < l; i++ {
		x := int32(s[i]) // permit sign extension to use full 64-bit range
		f.valEnc(buf, uint64(x))
	}

	o.Write(f.tagcode)
	o.EncodeVarint(uint64(buf.Len()))
	o.Write(buf.Bytes())
	return nil
}

// Encode a slice of uint32s ([]uint32) in packed format.
// Exactly the same as int32, except for no sign extension.
func (o *Buffer) enc_slice_uint32(f *Field, i interface{}) error {
	s, err := cast.ToIntSliceE(i)
	if err != nil {
		return errInvalidDataType
	}
	l := len(s)
	if l == 0 {
		return ErrNil
	}
	// TODO: Reuse a Buffer.
	buf := &Buffer{}
	for i := 0; i < l; i++ {
		x := int32(s[i])
		f.valEnc(buf, uint64(x))
	}

	o.Write(f.tagcode)
	o.EncodeVarint(uint64(buf.Len()))
	o.Write(buf.Bytes())
	return nil
}

// Encode a slice of int64s ([]int64) in packed format.
func (o *Buffer) enc_slice_int64(f *Field, i interface{}) error {
	s, err := ToInt64SliceE(i)
	if err != nil {
		return errInvalidDataType
	}
	l := len(s)
	if l == 0 {
		return ErrNil
	}
	// TODO: Reuse a Buffer.
	buf := &Buffer{}
	for i := 0; i < l; i++ {
		f.valEnc(buf, uint64(s[i]))
	}

	o.Write(f.tagcode)
	o.EncodeVarint(uint64(buf.Len()))
	o.Write(buf.Bytes())
	return nil
}

// Encode a slice of strings ([]string).
func (o *Buffer) enc_slice_string(f *Field, i interface{}) error {
	ss, err := cast.ToStringSliceE(i)
	if err != nil {
		return errInvalidDataType
	}
	l := len(ss)
	for i := 0; i < l; i++ {
		o.Write(f.tagcode)
		o.EncodeStringBytes(ss[i])
	}
	return nil
}
