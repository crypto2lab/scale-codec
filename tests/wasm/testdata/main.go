package main

// #include <stdlib.h>
import "C"

import (
	"bytes"
	"unsafe"

	scale_codec "github.com/crypto2lab/scale-codec"
)

func main() {}

func add(x, y *scale_codec.Integer[uint64]) uint64 {
	return x.Value + y.Value
}

//export addTwoNumbers
func addTwoNumbers(x, y uint64) (ptrSizeResult uint64) {
	fstPtr, fstSize := split(x)
	sndPtr, sndSize := split(y)

	fst := &scale_codec.Integer[uint64]{}
	snd := &scale_codec.Integer[uint64]{}

	err := fst.UnmarshalSCALE(bytes.NewReader(ptrToByteSlice(fstPtr, fstSize)))
	if err != nil {
		panic(err)
	}

	err = snd.UnmarshalSCALE(bytes.NewReader(ptrToByteSlice(sndPtr, sndSize)))
	if err != nil {
		panic(err)
	}

	result := &scale_codec.Integer[uint64]{Value: add(fst, snd)}
	encodedBytes, err := result.MarshalSCALE()
	if err != nil {
		panic(err)
	}

	ptr, size := sliceToPtrSize(encodedBytes)
	return join(ptr, size)
}

func split(ptrSize uint64) (ptr, size uint32) {
	return uint32(ptrSize >> 32), uint32(ptrSize & 0xFFFFFFFF)
}

func join(ptr, size uint32) (ptrSize uint64) {
	return (uint64(ptr) << uint64(32)) | uint64(size)
}

func ptrToByteSlice(ptr, size uint32) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), size)
}

func sliceToPtrSize(o []byte) (uint32, uint32) {
	size := C.ulong(len(o))
	ptr := unsafe.Pointer(C.malloc(size))
	copy(unsafe.Slice((*byte)(ptr), size), o)
	return uint32(uintptr(ptr)), uint32(size)
}
