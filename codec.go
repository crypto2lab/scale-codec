package scale_codec

import (
	"errors"
	"io"
)

var ErrUnexpectedOptionTag = errors.New("unexpected option tag")
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
