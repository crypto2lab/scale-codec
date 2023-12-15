package scale_codec

import (
	"bytes"
	"fmt"
	"io"
)

type Result struct {
	ok   Encodable
	isOk bool

	err   Encodable
	isErr bool
}

func NewResult(ok, err Encodable) *Result {
	return &Result{
		ok:  ok,
		err: err,
	}
}

func Ok(ok Encodable) *Result {
	return &Result{
		isOk: true,
		ok:   ok,
	}
}

func Err(err Encodable) *Result {
	return &Result{
		isErr: true,
		err:   err,
	}
}

func (r Result) MarshalSCALE() ([]byte, error) {
	if r.isErr {
		encErrResult := []byte{0x01}
		encErr, err := r.err.MarshalSCALE()
		if err != nil {
			return nil, fmt.Errorf("encoding result error: %w", err)
		}

		return bytes.Join([][]byte{encErrResult, encErr}, nil), nil
	}

	if r.isOk {
		encOkResult := []byte{0x00}
		encOk, err := r.ok.MarshalSCALE()
		if err != nil {
			return nil, fmt.Errorf("encoding result ok: %w", err)
		}

		return bytes.Join([][]byte{encOkResult, encOk}, nil), nil
	}

	return nil, ErrCannotEncodeEmptyResult
}

func (r *Result) UnmarshalSCALE(reader io.Reader) error {
	encResultTag := make([]byte, 1)
	n, err := reader.Read(encResultTag)
	if err != nil {
		return err
	}

	if n != 1 {
		return fmt.Errorf("%w: want %v, got %v", ErrUnexpectedReadBytes, 1, n)
	}

	var unmarshaler Encodable

	switch encResultTag[0] {
	case 0x00:
		*r = *Ok(r.ok)
		unmarshaler = r.ok
	case 0x01:
		*r = *Err(r.err)
		unmarshaler = r.err
	default:
		return fmt.Errorf("%w: %v", ErrUnexpectedResultTag, encResultTag)
	}

	return unmarshaler.UnmarshalSCALE(reader)
}

func (r *Result) Unwrap() Encodable {
	if r.isErr {
		panic(r.err)
	}

	if r.isOk {
		return r.ok
	}

	panic("cannot unwrap empty result")
}

func (r *Result) IsErr() bool {
	return r.isErr
}

func (r *Result) Err() Encodable {
	return r.err
}
