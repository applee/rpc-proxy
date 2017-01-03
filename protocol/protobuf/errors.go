package protobuf

import (
	"errors"
	"fmt"
)

// RequiredNotSetError is the error returned if Marshal is called with
// a protocol buffer struct whose required fields have not
// all been initialized. It is also the error returned if Unmarshal is
// called with an encoded protocol buffer that does not include all the
// required fields.
//
// When printed, RequiredNotSetError reports the first unset required field in a
// message. If the field cannot be precisely determined, it is reported as
// "{Unknown}".
type RequiredNotSetError struct {
	field string
}

func (e *RequiredNotSetError) Error() string {
	return fmt.Sprintf("proto: required field %q not set", e.field)
}

type MessageNotExistsError struct {
	name string
}

func (e *MessageNotExistsError) Error() string {
	return fmt.Sprintf("proto: no message named %s", e.name)
}

var (
	// errInvalidParser is the error returned if the given name parser do not
	// exists.
	errInvalidParser = errors.New("invalid parser")

	// errRepeatedHasNil is the error returned if Marshal is called with
	// a struct with a repeated field containing a nil element.
	errRepeatedHasNil = errors.New("proto: repeated field has nil element")

	// errInvalidField is the error returned if the field's type of message is
	// an unknown type.
	errInvalidField = errors.New("proto: invalid field type")

	// errInvalidDataType is the error returned if the data of message is not match
	// with the proto
	errInvalidDataType = errors.New("proto: invalid data to marshal")

	// errOneofHasNil is the error returned if Marshal is called with
	// a struct with a oneof field containing a nil element.
	errOneofHasNil = errors.New("proto: oneof field has nil value")

	// ErrNil is the error returned if Marshal is called with nil.
	ErrNil = errors.New("proto: Marshal called with nil")

	// ErrTooLarge is the error returned if Marshal is called with a
	// message that encodes to >2GB.
	ErrTooLarge = errors.New("proto: message encodes to over 2 GB")
)
