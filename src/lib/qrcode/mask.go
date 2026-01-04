package qrcode

import (
	"math"
	"yotudo/src/lib/qrcode/bitbuffer"
)

type maskFunc func(r, c int) bool

var maskFuncs = [maskCount]maskFunc{
	func(x, y int) bool { return (y+x)%2 == 0 },
	func(x, y int) bool { return y%2 == 0 },
	func(x, y int) bool { return x%3 == 0 },
	func(x, y int) bool { return (y+x)%3 == 0 },
	func(x, y int) bool { return (y/2+x/3)%2 == 0 },
	func(x, y int) bool { return (y*x)%2+(y*x)%3 == 0 },
	func(x, y int) bool { return ((y*x)%2+(y*x)%3)%2 == 0 },
	func(x, y int) bool { return ((y+x)%2+(y*x)%3)%2 == 0 },
}

// consecutive modules
func penaltyN1(m *QR) int {
	score := 0

	// Rows
	for y := range m.Size {
		run := 1

		for x := 1; x < m.Size; x++ {
			if m.Get(x, y).IsSet() == m.Get(x-1, y).IsSet() {
				run++

				if run == 5 {
					score += 3
				} else if run > 5 {
					score++
				}
			} else {
				run = 1
			}
		}
	}

	// Columns
	for x := range m.Size {
		run := 1

		for y := 1; y < m.Size; y++ {
			if m.Get(x, y).IsSet() == m.Get(x, y-1).IsSet() {
				run++
				if run == 5 {
					score += 3
				} else if run > 5 {
					score++
				}
			} else {
				run = 1
			}
		}
	}
	return score
}

// 2x2 blocks
func penaltyN2(m *QR) int {
	score := 0

	for y := range m.Size - 1 {
		for x := range m.Size - 1 {
			v := m.Get(x, y).IsSet()

			if m.Get(x+1, y).IsSet() == v &&
				m.Get(x, y+1).IsSet() == v &&
				m.Get(x+1, y+1).IsSet() == v {
				score += 3
			}
		}
	}

	return score
}

// finder-like patterns
func penaltyN3(m *QR) int {
	score := 0
	pattern1 := []bool{true, false, true, true, true, false, true, false, false, false, false}
	pattern2 := []bool{false, false, false, false, true, false, true, true, true, false, true}

	for y := range m.Size {
		for x := range m.Size - 12 {
			match := true

			for i := range 11 {
				if m.Get(x+i, y).IsSet() != pattern1[i] &&
					m.Get(x+i, y).IsSet() != pattern2[i] {
					match = false
					break
				}
			}

			if match {
				score += 40
			}
		}
	}

	for x := range m.Size {
		for y := range m.Size - 12 {
			match := true

			for i := range 11 {
				if m.Get(x, y+i).IsSet() != pattern1[i] &&
					m.Get(x, y+i).IsSet() != pattern2[i] {
					match = false
					break
				}
			}
			if match {
				score += 40
			}
		}
	}
	return score
}

// balance of black/white
func penaltyN4(m *QR) int {
	black := 0
	total := float64(m.Size * m.Size)

	for x := 0; x < m.Size; x++ {
		for y := 0; y < m.Size; y++ {
			if m.Get(x, y).IsSet() {
				black++
			}
		}
	}

	percent := float64(black*100) / total

	return int(math.Abs(percent-50)) / 5 * 10
}

func totalPenalty(m *QR) int {
	return penaltyN1(m) +
		penaltyN2(m) +
		penaltyN3(m) +
		penaltyN4(m)
}

func formatBits(mask int) int {
	const lowRecoveryLevelFormatId int = 0x8
	data := lowRecoveryLevelFormatId | mask
	bch := bchEncode(data, 0b10100110111, 10)
	return ((data << 10) | bch) ^ 0b101010000010010
}

func bchEncode(data int, poly int, polyLen int) int {
	data <<= polyLen
	for bitLen(data) >= polyLen+1 {
		shift := bitLen(data) - polyLen - 1
		data ^= poly << shift
	}
	return data
}

func bitLen(v int) int {
	n := 0
	for v > 0 {
		v >>= 1
		n++
	}
	return n
}

func createQRWithMask(buffer *bitbuffer.BitBuffer, metadata *metadata, mask int) (*QR, error) {
	dst := newQR(metadata)
	maskFunc := maskFuncs[mask]

	dx := 1
	dy := -1
	x := int(dst.Size - 2)
	y := int(dst.Size - 1)

	safeguardLimit := dst.Size * dst.Size / 2

	for i := uint(0); i < buffer.Len(); i++ {
		mValue := maskFunc(x+dx, y)
		bitValue, err := buffer.At(i)
		if err != nil {
			return nil, err
		}

		value := mValue != bitValue
		dst.Get(x+dx, y).SetCell(value)
		dst.Get(x+dx, y).SetReserved(value)

		if i == buffer.Len()-1 {
			break
		}

		for range safeguardLimit {
			if dx == 0 {
				if (y <= 0 && dy < 0) || (y >= dst.Size-1 && dy > 0) {
					dy = -dy
					x -= 2
				} else {
					y += dy
				}
			}

			dx = (dx + 1) & 1

			if x == 5 {
				x--
			}

			if !dst.Get(x+dx, y).IsReserved() {
				break
			}
		}
	}

	return dst, nil
}
