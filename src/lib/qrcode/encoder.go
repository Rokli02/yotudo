package qrcode

import (
	"fmt"
	"yotudo/src/lib/qrcode/bitbuffer"
)

func encodeDataIntoBitBuffer(data []byte, metadata *metadata) (*bitbuffer.BitBuffer, error) {
	if len(data) <= 0 {
		return nil, fmt.Errorf("input data is empty")
	}

	buffer := bitbuffer.New()
	buffer.AppendByte(metadata.Mode.indicator(), 4)

	ccBits := charCountBitLength(metadata.Mode, metadata.Version)
	buffer.AppendUInt32(metadata.CharCount, ccBits)

	switch metadata.Mode {
	case ModeNumeric:
		for i := 0; i < len(data); i += 3 {
			charsRemaining := len(data) - i

			var value uint32
			var bitsUsed uint = 1

			for j := 0; j < charsRemaining && j < 3; j++ {
				value *= 10
				value += uint32(data[i+j] - 0x30)
				bitsUsed += 3
			}

			buffer.AppendUInt32(value, bitsUsed)
		}
	case ModeAlphanumeric:
		for i := 0; i < len(data); i += 2 {
			charsRemaining := len(data) - i
			if charsRemaining > 1 {
				v1 := encodeAlphanumericCharacter(data[i])
				v2 := encodeAlphanumericCharacter(data[i+1])

				buffer.AppendUInt32(v1*45+v2, 11)
			} else {
				v := encodeAlphanumericCharacter(data[i])

				buffer.AppendUInt32(v, 6)
			}
		}
	case ModeByte:
		if err := buffer.AppendBytes(data); err != nil {
			return nil, err
		}
	}

	// Terminator
	capacity := dataBitsL[metadata.Version-1]
	remaining := capacity - buffer.Len()
	if remaining > 0 {
		term := min(4, remaining)
		if err := buffer.AppendNBits(term, false); err != nil {
			return nil, err
		}
	}

	// Pad to byte boundary
	if setBits := buffer.Len() % 8; setBits != 0 {
		if err := buffer.AppendNBits(8-setBits, false); err != nil {
			return nil, err
		}
	}

	// Pad bytes
	for i := 0; buffer.Len() < capacity; i++ {
		if err := buffer.AppendByte(padBytes[i%2], 8); err != nil {
			return nil, err
		}
	}

	return buffer, nil
}

var padBytes = []byte{0xEC, 0x11}

var dataBitsL = []uint{
	152, 272, 440, 640, 864, 1088, 1248, 1552, 1856, 2192,
	2592, 2960, 3424, 3688, 4184, 4712, 5176, 5768, 6360, 6888,
	7456, 8048, 8752, 9392, 10208, 10960, 11744, 12248, 13048, 13880,
	14744, 15640, 16568, 17528, 18448, 19472, 20528, 21616, 22496, 23648,
}

func charCountBitLength(mode encodingMode, version version) uint {
	switch mode {
	case ModeNumeric:
		if version <= 9 {
			return 10
		}

		if version <= 26 {
			return 12
		}

		return 14
	case ModeAlphanumeric:
		if version <= 9 {
			return 9
		}

		if version <= 26 {
			return 11
		}

		return 13
	default:
		if version <= 9 {
			return 8
		}

		return 16
	}
}

func encodeAlphanumericCharacter(v byte) uint32 {
	c := uint32(v)

	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'A' && c <= 'Z':
		return c - 'A' + 10
	case c == ' ':
		return 36
	case c == '$':
		return 37
	case c == '%':
		return 38
	case c == '*':
		return 39
	case c == '+':
		return 40
	case c == '-':
		return 41
	case c == '.':
		return 42
	case c == '/':
		return 43
	case c == ':':
		return 44
	default:
		fmt.Printf("encodeAlphanumericCharacter() with non alphanumeric char %v.", v)
	}

	return 0
}
