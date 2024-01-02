package scale_codec

import (
	"bytes"
	"fmt"
	"io"
)

type ResultG[T Marshaler, E Marshaler] struct {
	ok   T
	isOk bool

	err   E
	isErr bool
}

func UnmarshalResultFromRawBytes[T Marshaler, E Marshaler](
	okF func(io.Reader) (T, error),
	errF func(io.Reader) (E, error)) func(reader io.Reader) (*ResultG[T, E], error) {
	return func(reader io.Reader) (*ResultG[T, E], error) {
		result := &ResultG[T, E]{}
		err := result.UnmarshalSCALE(reader, okF, errF)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func (r ResultG[T, E]) MarshalSCALE() ([]byte, error) {
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

func (r *ResultG[T, E]) UnmarshalSCALE(reader io.Reader,
	okF func(io.Reader) (T, error), errF func(io.Reader) (E, error)) error {
	encResultTag := make([]byte, 1)
	n, err := reader.Read(encResultTag)
	if err != nil {
		return err
	}

	if n != 1 {
		return fmt.Errorf("%w: want %v, got %v", ErrUnexpectedReadBytes, 1, n)
	}

	switch encResultTag[0] {
	case 0x00:
		ok, err := okF(reader)
		if err != nil {
			return fmt.Errorf("while parsing result ok: %w", err)
		}
		*r = *OkG[T, E](ok)
	case 0x01:
		errResult, err := errF(reader)
		if err != nil {
			return fmt.Errorf("while parsing result err: %w", err)
		}
		*r = *ErrG[T, E](errResult)
		return nil
	default:
		return fmt.Errorf("%w: %v", ErrUnexpectedResultTag, encResultTag)
	}

	return nil
}

func OkG[T Marshaler, E Marshaler](ok T) *ResultG[T, E] {
	return &ResultG[T, E]{
		isOk: true,
		ok:   ok,
	}
}

func ErrG[T Marshaler, E Marshaler](err E) *ResultG[T, E] {
	return &ResultG[T, E]{
		isErr: true,
		err:   err,
	}
}

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
