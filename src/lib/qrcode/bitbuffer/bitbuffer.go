package bitbuffer

import (
	"fmt"
	"strings"
)

const (
	byte_mask byte = 0x80
)

type BitBuffer struct {
	data     []byte
	bitsUsed uint
}

func New() *BitBuffer {
	b := new(BitBuffer)

	b.data = make([]byte, 0)

	return b
}

func Clone(b *BitBuffer) *BitBuffer {
	return &BitBuffer{bitsUsed: b.bitsUsed, data: b.data[:]}
}

func (b *BitBuffer) ensureBufferLength(appendBits uint) {
	bitsToUse := b.bitsUsed + appendBits
	minDataLength := bitsToUse / 8

	if bitsToUse%8 != 0 {
		minDataLength++
	}

	if minDataLength > uint(len(b.data)) {
		appendSize := max(2*len(b.data), int(minDataLength))
		b.data = append(b.data, make([]byte, appendSize)...)
	}

}

func (b *BitBuffer) AppendUInt32(value uint32, numBits uint) error {
	const MAX_BIT_SIZE = 32

	if numBits > MAX_BIT_SIZE {
		return fmt.Errorf("numBits in 'AppendUInt32' is too large (numBits=%d)", numBits)
	}

	b.ensureBufferLength(numBits)

	for i := int(numBits) - 1; i >= 0; i-- {
		eee := value & (1 << i)
		if eee != 0 {
			dataIndex := b.bitsUsed / 8
			bitIndex := b.bitsUsed % 8

			b.data[dataIndex] |= byte_mask >> bitIndex
		}

		b.bitsUsed++
	}

	return nil
}

func (b *BitBuffer) AppendByte(value byte, numBits uint) error {
	const MAX_BIT_SIZE = 8

	if numBits > MAX_BIT_SIZE {
		return fmt.Errorf("numBits in 'AppendByte' is too large (numBits=%d)", numBits)
	}

	b.ensureBufferLength(numBits)

	for i := int(numBits) - 1; i >= 0; i-- {
		if value&(1<<i) != 0 {
			dataIndex := b.bitsUsed / 8
			bitIndex := b.bitsUsed % 8

			b.data[dataIndex] |= byte_mask >> bitIndex
		}

		b.bitsUsed++
	}

	return nil
}

func (b *BitBuffer) AppendBytes(values []byte) error {
	for _, value := range values {
		if err := b.AppendByte(value, 8); err != nil {
			return err
		}
	}

	return nil
}

func (b *BitBuffer) AppendNBits(N uint, bitValue bool) error {
	b.ensureBufferLength(N)

	if !bitValue {
		b.bitsUsed += N

		return nil
	}

	for range N {
		dataIndex := b.bitsUsed / 8
		bitIndex := b.bitsUsed % 8
		b.data[dataIndex] |= byte_mask >> bitIndex

		b.bitsUsed++
	}

	return nil
}

func (b *BitBuffer) At(index uint) (bool, error) {
	if index >= b.bitsUsed {
		return false, fmt.Errorf("index %d out of range (max=%d)", index, b.bitsUsed-1)
	}

	return (b.data[index/8] & (0x80 >> byte(index%8))) != 0, nil
}

func (b *BitBuffer) ByteAt(index uint) byte {
	var result byte

	for i := index; i < index+8 && i < b.bitsUsed; i++ {
		result <<= 1
		if bit, _ := b.At(i); bit {
			result |= 1
		}
	}

	return result
}

func (b *BitBuffer) Len() uint {
	return b.bitsUsed
}

func (b *BitBuffer) ByteLen() uint {
	length := b.bitsUsed / 8

	if b.bitsUsed%8 != 0 {
		length++
	}

	return length
}

func (b *BitBuffer) String() string {
	sb := strings.Builder{}

	for i, value := range b.data {
		activeBitsInByte := int(b.bitsUsed) - i*8

		if activeBitsInByte <= 0 {
			break
		}

		for i := range min(8, activeBitsInByte) {
			if value&(0x80>>i) != 0 {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('0')
			}
		}
	}

	return sb.String()
}

func (b *BitBuffer) GetBitChunk(from, size uint) ([]byte, uint, error) {
	if from >= b.bitsUsed {
		return nil, 0, fmt.Errorf("'to' is out of bound for \"SliceBytes\" (%d/%d)", from, len(b.data))
	}

	chunkLength := size / 8
	if size%8 != 0 {
		chunkLength++
	}
	chunk := make([]byte, chunkLength)
	var actualSize uint

	for i := uint(from); i < from+size; i++ {
		bitValue, err := b.At(i)
		if err != nil {
			return chunk, actualSize, nil
		}

		if bitValue {
			bitIndex := actualSize % 8
			chunk[actualSize/8] |= 0x80 >> bitIndex
		}

		actualSize++
	}

	return chunk, size, nil
}

func (b *BitBuffer) SliceBytes(fromByte, toByte uint) (*BitBuffer, error) {
	if toByte < fromByte {
		return nil, fmt.Errorf("'toByte' must be greater than 'fromByte' in \"SliceBytes\" (fromByte=%d, toByte=%d)", fromByte, toByte)
	}

	if toByte > uint(len(b.data)) {
		return nil, fmt.Errorf("'toByte' out of bound for \"SliceBytes\" (%d/%d)", toByte, len(b.data))
	}

	newBitBuffer := New()
	newBitBuffer.bitsUsed = (toByte - fromByte) * 8
	newBitBuffer.data = b.data[fromByte:toByte]

	return newBitBuffer, nil
}

func (b *BitBuffer) Data() []bool {
	data := make([]bool, b.bitsUsed)

	for i := range b.bitsUsed {
		dataIndex := i / 8
		bitIndex := i % 8

		data[i] = (b.data[dataIndex] & (byte_mask >> bitIndex)) != 0
	}

	return data
}
