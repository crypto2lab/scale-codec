package scale_codec

import (
	"errors"
	"io"
)

var ErrUnexpectedResultTag = errors.New("unexpected result tag")
var ErrCannotEncodeEmptyResult = errors.New("cannot encode empty result")
var ErrCannotEncodeEmptyOption = errors.New("cannot encode empty option")
var ErrUnexpectedReadBytes = errors.New("unexpected read bytes")

type Marshaler interface {
	MarshalSCALE() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalSCALE(io.Reader) error
}

type Encodable interface {
	Marshaler
	Unmarshaler
}

type Option[T Encodable] struct {
	inner  T
	isNone bool
}

type Vector[T Encodable] struct {
	items []T
}

type String struct {
	inner string
}

type Tuple[T Encodable] struct {
	items []T
	size  uint
}

type Struct struct{}

type Enum struct{}