package scale_codec

import (
	"fmt"
	"io"
)

type FieldAccess interface {
	Field(int) Encodable
}

type Tuple struct {
	Items []Encodable
}

func NewTuple(items ...Encodable) *Tuple {
	return &Tuple{items}
}

func (t Tuple) MarshalSCALE() ([]byte, error) {
	encodedTuple := make([]byte, 0)
	for idx, item := range t.Items {
		encodedItem, err := item.MarshalSCALE()
		if err != nil {
			return nil, fmt.Errorf("encoding item at index %v: %w", idx, err)
		}

		encodedTuple = append(encodedTuple, encodedItem...)
	}

	return encodedTuple, nil
}

func (t *Tuple) UnmarshalSCALE(reader io.Reader) error {
	for idx, item := range t.Items {
		err := item.UnmarshalSCALE(reader)
		if err != nil {
			return fmt.Errorf("encoding item at index %v: %w", idx, err)
		}
	}

	return nil
}

func (t *Tuple) FieldAccess(at int) Encodable {
	if at > len(t.Items)-1 {
		panic(fmt.Sprintf("cannot access tuple item at: %d", at))
	}

	return t.Items[at]
}
