package messages

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestNewMessageArgInt32(t *testing.T) {
	number := int32(rand.Int())
	byteRepr := make([]byte, 0)
	byteRepr = append(byteRepr, Int32Type)
	byteRepr = append(byteRepr, byte(number&0xFF))
	byteRepr = append(byteRepr, byte((number>>8)&0xFF))
	byteRepr = append(byteRepr, byte((number>>16)&0xFF))
	byteRepr = append(byteRepr, byte((number>>24)&0xFF))

	arg := NewMessageArg(number)
	if !bytes.Equal(arg.Body, byteRepr) {
		t.Errorf("Result mismatch: expected %v, got %v", byteRepr, arg.Body)
	}
}

func TestNewMessageArgInt64(t *testing.T) {
	number := int64(rand.Uint64())

	byteRepr := make([]byte, 0)
	byteRepr = append(byteRepr, Int64Type)
	byteRepr = append(byteRepr, byte(number&0xFF))
	byteRepr = append(byteRepr, byte((number>>8)&0xFF))
	byteRepr = append(byteRepr, byte((number>>16)&0xFF))
	byteRepr = append(byteRepr, byte((number>>24)&0xFF))
	byteRepr = append(byteRepr, byte((number>>32)&0xFF))
	byteRepr = append(byteRepr, byte((number>>40)&0xFF))
	byteRepr = append(byteRepr, byte((number>>48)&0xFF))
	byteRepr = append(byteRepr, byte((number>>56)&0xFF))

	arg := NewMessageArg(number)
	if !bytes.Equal(arg.Body, byteRepr) {
		t.Errorf("Result mismatch: expected %v, got %v", byteRepr, arg.Body)
	}
}

func TestNewMessageArgString(t *testing.T) {
	str := []byte("Lorem ipsum dolor sit amet")
	res := make([]byte, 0)
	res = append(res, StringType)
	res = append(res, str...)
	res = append(res, 0)

	arg := NewMessageArg(str)
	if !bytes.Equal(arg.Body, res) {
		t.Errorf("Result mismatch: expected %v, got %v", res, arg.Body)
	}
}

func TestNewMessageArgDefault(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Errorf("No panic occured")
		}
	}()

	_ = NewMessageArg(complex64(1))
}
