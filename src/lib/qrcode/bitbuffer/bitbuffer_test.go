package bitbuffer_test

import (
	"testing"
	"yotudo/src/lib/qrcode/bitbuffer"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}

func noError(err error) {
	if err != nil {
		panic(err)
	}
}

func TestBitBufferAt(t *testing.T) {
	bb := bitbuffer.New()

	noError(bb.AppendByte(0b0100_1010, 8))

	b1 := must(bb.At(1))
	b2 := must(bb.At(4))
	b3 := must(bb.At(5))

	if b1 && !b2 && b3 {
		t.Errorf("BitBuffer.At() misbehaves - %t-%t-%t", b1, b2, b3)
		return
	}
}

func TestBitBufferString(t *testing.T) {
	bb := bitbuffer.New()

	noError(bb.AppendByte(0b0100_1010, 8))

	stringifiedBigBuffer := bb.String()

	if stringifiedBigBuffer != "01001010" {
		t.Errorf("BitBuffer.String() misbehaves (\"%s\")", stringifiedBigBuffer)
		return
	}
}

func TestBitBufferString_2(t *testing.T) {
	bb := bitbuffer.New()

	noError(bb.AppendByte(0b0100_1010, 4))

	stringifiedBigBuffer := bb.String()

	if stringifiedBigBuffer != "1010" {
		t.Errorf("BitBuffer.String() misbehaves (\"%s\")", stringifiedBigBuffer)
		return
	}
}

func TestBitBufferString_3(t *testing.T) {
	bb := bitbuffer.New()

	noError(bb.AppendByte(0b0100_1010, 6))

	stringifiedBigBuffer := bb.String()

	if stringifiedBigBuffer != "001010" {
		t.Errorf("BitBuffer.String() misbehaves (\"%s\")", stringifiedBigBuffer)
		return
	}
}

func TestBitBufferAppendByte(t *testing.T) {
	bb := bitbuffer.New()

	noError(bb.AppendByte(0b0100_1010, 6))
	noError(bb.AppendByte(0b0100_1010, 4))

	if bb.Len() != 10 {
		t.Errorf("BitBuffer.Len() expected to be %d, but was %d", 10, bb.Len())
		return
	}

	stringifiedBigBuffer := bb.String()

	if stringifiedBigBuffer != "0010101010" {
		t.Errorf("BitBuffer.String() misbehaves (\"%s\")", stringifiedBigBuffer)
		return
	}
}

func TestBitBufferAppendByte_Fail(t *testing.T) {
	bb := bitbuffer.New()

	err := bb.AppendByte(0b0100_1010, 9)

	if err == nil {
		t.Error("BitBuffer.AppendByte() should have failed")
		return
	}
}

func TestBitBufferAppendBytes(t *testing.T) {
	bb := bitbuffer.New()
	var appendBits uint = 16

	noError(bb.AppendBytes([]byte{0b0100_1010, 0b0001_1011}))

	if bb.Len() != appendBits {
		t.Errorf("BitBuffer.Len() expected to be %d, but was %d", appendBits, bb.Len())
		return
	}

	stringifiedBigBuffer := bb.String()

	if stringifiedBigBuffer != "0100101000011011" {
		t.Errorf("BitBuffer.String() misbehaves (\"%s\")", stringifiedBigBuffer)
		return
	}
}

func TestBitBufferAppendUInt32(t *testing.T) {
	bb := bitbuffer.New()
	var appendBits uint = 11

	noError(bb.AppendUInt32(0b1110010, appendBits))

	if bb.Len() != appendBits {
		t.Errorf("BitBuffer.Len() expected to be %d, but was %d", appendBits, bb.Len())
		return
	}

	stringifiedBigBuffer := bb.String()

	if stringifiedBigBuffer != "00001110010" {
		t.Errorf("BitBuffer.String() misbehaves (\"%s\")", stringifiedBigBuffer)
		return
	}
}

func TestBitBufferData(t *testing.T) {
	bb := bitbuffer.New()
	var appendBits uint = 11
	expectedData := []bool{
		false, false, false, false,
		true, true, true, false,
		false, true, false,
	}

	noError(bb.AppendUInt32(0b1110010, appendBits))

	result := bb.Data()

	for i := range result {
		if result[i] != expectedData[i] {
			t.Error("BitBuffer.Data() didn't return with expected slice")
			return
		}
	}
}

func TestBitBufferAppendNBits(t *testing.T) {
	bb := bitbuffer.New()
	var appendBits uint = 11
	expectedData := []bool{
		false, false, false, false,
		true, true, true, false,
		false, true, false,
		false, false, false, false, false,
		false, false, false, false,
	}

	noError(bb.AppendUInt32(0b1110010, appendBits))
	noError(bb.AppendNBits(9, false))

	appendBits += 9

	if bb.Len() != appendBits {
		t.Errorf("BitBuffer.Len() expected to be %d, but was %d", appendBits, bb.Len())
		return
	}

	result := bb.Data()

	for i := range result {
		if result[i] != expectedData[i] {
			t.Error("BitBuffer.Data() didn't return with expected slice")
			return
		}
	}
}

func TestBitBufferAppendNBits_2(t *testing.T) {
	bb := bitbuffer.New()
	var appendBits uint = 11
	expectedData := []bool{
		false, false, false, false,
		true, true, true, false,
		false, true, false,
		true, true, true, true, true,
		true, true, true, true,
	}

	noError(bb.AppendUInt32(0b1110010, appendBits))
	noError(bb.AppendNBits(9, true))

	appendBits += 9

	if bb.Len() != appendBits {
		t.Errorf("BitBuffer.Len() expected to be %d, but was %d", appendBits, bb.Len())
		return
	}

	result := bb.Data()

	for i := range result {
		if result[i] != expectedData[i] {
			t.Error("BitBuffer.Data() didn't return with expected slice")
			return
		}
	}
}

func TestBitBufferGetBitChunk(t *testing.T) {
	bb := bitbuffer.New()
	var appendBits uint = 11
	var chunkToReturn byte = 0b00001110

	noError(bb.AppendUInt32(0b00001110010, appendBits))

	chunk, read, err := bb.GetBitChunk(0, 8)
	if err != nil {
		t.Errorf("BitBuffer.GetBitChunk() misbehaves (\"%v\")", err)
		return
	}

	if read != 8 {
		t.Errorf("BitBuffer.GetBitChunk() misbehaves, didn't read enough (\"%d\")", read)
		return
	}

	if len(chunk) < 1 {
		t.Errorf("BitBuffer.GetBitChunk() misbehaves, no chunk returned (\"%d\")", read)
		return
	}

	if chunk[0] != chunkToReturn {
		t.Errorf("BitBuffer.GetBitChunk() misbehaves, wrong chunk returned expected '%08b', but got '%08b'", chunkToReturn, chunk[0])
		return
	}
}

func TestBitBufferGetBitChunk_ReadLessThanExpected(t *testing.T) {
	bb := bitbuffer.New()
	var appendBits uint = 11
	var chunkToReturn byte = 0b10010000

	noError(bb.AppendUInt32(0b00001110010, appendBits))

	chunk, read, err := bb.GetBitChunk(6, 8)
	if err != nil {
		t.Errorf("BitBuffer.GetBitChunk() misbehaves (\"%v\")", err)
		return
	}

	if read == 8 {
		t.Errorf("BitBuffer.GetBitChunk() misbehaves, read too much (\"%d\")", read)
		return
	}

	if len(chunk) < 1 {
		t.Errorf("BitBuffer.GetBitChunk() misbehaves, no chunk returned (\"%d\")", read)
		return
	}

	if chunk[0] != chunkToReturn {
		t.Errorf("BitBuffer.GetBitChunk() misbehaves, wrong chunk returned expected '%08b', but got '%08b'", chunkToReturn, chunk[0])
		return
	}
}
