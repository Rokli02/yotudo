package qrcode

import (
	"fmt"
	"yotudo/src/lib/qrcode/bitbuffer"
)

type Cell uint8

const (
	set_mask Cell = 1 << iota
	reserved_mask
)

func (c *Cell) SetReserved(v bool) {
	if v {
		*c |= reserved_mask
	} else {
		*c &= ^reserved_mask
	}
}

func (c *Cell) SetCell(v bool) {
	if v {
		*c |= set_mask
	} else {
		*c &= ^set_mask
	}
}

func (c Cell) IsSet() bool {
	return c&set_mask != 0
}

func (c Cell) IsReserved() bool {
	return c&reserved_mask != 0
}

func (c *QR) AttributesJson() string {
	return fmt.Sprintf("{ value: %d, reserved: %d }", set_mask, reserved_mask)
}

func (qr *QR) Get(x, y int) *Cell {
	return &qr.Data[x+qr.Size*y]
}

const (
	finderPatternSize = 7
	maskCount         = 8
)

var (
	alignmentPositionsTable = [][]int{
		nil,
		{6, 18},
		{6, 22},
		{6, 26},
		{6, 30},
		{6, 34},
		{6, 22, 38},
		{6, 24, 42},
		{6, 26, 46},
		{6, 28, 50},
		{6, 30, 54},
		{6, 32, 58},
		{6, 34, 62},
		{6, 26, 46, 66},
		{6, 26, 48, 70},
		{6, 26, 50, 74},
		{6, 30, 54, 78},
		{6, 30, 56, 82},
		{6, 30, 58, 86},
		{6, 34, 62, 90},
		{6, 28, 50, 72, 94},
		{6, 26, 50, 74, 98},
		{6, 30, 54, 78, 102},
		{6, 28, 54, 80, 106},
		{6, 32, 58, 84, 110},
		{6, 30, 58, 86, 114},
		{6, 34, 62, 90, 118},
		{6, 26, 50, 74, 98, 122},
		{6, 30, 54, 78, 102, 126},
		{6, 26, 52, 78, 104, 130},
		{6, 30, 56, 82, 108, 134},
		{6, 34, 60, 86, 112, 138},
		{6, 30, 58, 86, 114, 142},
		{6, 34, 62, 90, 118, 146},
		{6, 30, 54, 78, 102, 126, 150},
		{6, 24, 50, 76, 102, 128, 154},
		{6, 28, 54, 80, 106, 132, 158},
		{6, 32, 58, 84, 110, 136, 162},
		{6, 26, 54, 82, 110, 138, 166},
		{6, 30, 58, 86, 114, 142, 170},
	}

	versionBitSequence = []uint32{
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0x7c94,
		0x85bc,
		0x9a99,
		0xa4d3,
		0xbbf6,
		0xc762,
		0xd847,
		0xe60d,
		0xf928,
		0x10b78,
		0x1145d,
		0x12a17,
		0x13532,
		0x149a6,
		0x15683,
		0x168c9,
		0x177ec,
		0x18ec4,
		0x191e1,
		0x1afab,
		0x1b08e,
		0x1cc1a,
		0x1d33f,
		0x1ed75,
		0x1f250,
		0x209d5,
		0x216f0,
		0x228ba,
		0x2379f,
		0x24b0b,
		0x2542e,
		0x26a64,
		0x27541,
		0x28c69,
	}
)

func createQRWithBestMask(buffer *bitbuffer.BitBuffer, metadata *metadata) (*QR, int) {
	bestScore := 1<<31 - 1
	bestMask := 0
	var bestMatrix *QR

	for mask := range maskCount {
		candidate, err := createQRWithMask(buffer, metadata, mask)
		if err != nil {
			// TODO: Valami spécibb lekezelés
			fmt.Printf("Unable to create QR with mask number %d\n", mask)
			fmt.Println(err)

			continue
		}

		candidate.writeFormatInfo(mask)
		score := totalPenalty(candidate)

		if score < bestScore {
			bestScore = score
			bestMask = mask
			bestMatrix = candidate
		}
	}

	return bestMatrix, bestMask
}

func newQR(metadata *metadata) *QR {
	size := metadata.size()

	qr := new(QR)
	qr.Size = int(size)
	qr.Data = make([]Cell, size*size)

	qr.placeFinders()
	qr.placeTiming()
	qr.placeAlignment(metadata.Version)
	qr.reserveFormat()
	qr.reserveVersion(metadata.Version)

	return qr
}

func (qr *QR) placeFinder(x, y int) {
	for dy := -1; dy <= finderPatternSize; dy++ {
		for dx := -1; dx <= finderPatternSize; dx++ {
			ix, iy := x+dx, y+dy
			if iy < 0 || iy >= qr.Size || ix < 0 || ix >= qr.Size {
				continue
			}

			qr.Get(ix, iy).SetReserved(true)

			if dx == -1 || dx == finderPatternSize || dy == -1 || dy == finderPatternSize {
				qr.Get(ix, iy).SetCell(false)
			} else if dx == 0 || dx == 6 || dy == 0 || dy == 6 {
				qr.Get(ix, iy).SetCell(true)
			} else if dx >= 2 && dx <= 4 && dy >= 2 && dy <= 4 {
				qr.Get(ix, iy).SetCell(true)
			} else {
				qr.Get(ix, iy).SetCell(false)
			}
		}
	}
}

func (qr *QR) placeFinders() {
	qr.placeFinder(0, 0)
	qr.placeFinder(qr.Size-7, 0)
	qr.placeFinder(0, qr.Size-7)
}

func (qr *QR) placeTiming() {
	for i := 8; i < qr.Size-8; i++ {
		val := i%2 != 1

		qr.Get(6, i).SetCell(val)
		qr.Get(6, i).SetReserved(true)

		qr.Get(i, 6).SetCell(val)
		qr.Get(i, 6).SetReserved(true)
	}
}

func (qr *QR) placeAlignment(version version) {
	pos := alignmentPositionsTable[version-1]
	if len(pos) == 0 {
		return
	}

	for _, y := range pos {
		for _, x := range pos {
			// Skip overlapping finder areas
			if qr.Get(x, y).IsReserved() {
				continue
			}
			qr.placeAlignmentPattern(x, y)
		}
	}
}

func (qr *QR) placeAlignmentPattern(x, y int) {
	for dy := -2; dy <= 2; dy++ {
		for dx := -2; dx <= 2; dx++ {
			ix, iy := x+dx, y+dy
			qr.Get(ix, iy).SetReserved(true)

			if dx == -2 || dx == 2 || dy == -2 || dy == 2 {
				qr.Get(ix, iy).SetCell(true)
			} else if dx == 0 && dy == 0 {
				qr.Get(ix, iy).SetCell(true)
			} else {
				qr.Get(ix, iy).SetCell(false)
			}
		}
	}
}

func (qr *QR) reserveFormat() {
	// Format info around top-left finder
	for i := range 9 {
		qr.Get(8, i).SetReserved(true)
		qr.Get(i, 8).SetReserved(true)
	}

	// Top-right format
	for i := range 8 {
		qr.Get(qr.Size-1-i, finderPatternSize+1).SetReserved(true)
	}

	// Bottom-left format
	for i := range 7 {
		qr.Get(finderPatternSize+1, qr.Size-1-i).SetReserved(true)
	}

	qr.Get(finderPatternSize+1, qr.Size-finderPatternSize-1).SetReserved(true)
	qr.Get(finderPatternSize+1, qr.Size-finderPatternSize-1).SetCell(true)
}

func (qr *QR) reserveVersion(version version) {
	if version < 7 {
		return
	}

	vbits := versionBitSequence[version]

	for i := range 6 {
		for j := range 3 {
			bit := ((vbits >> (i*3 + j)) & 1) != 0
			qr.Get(i, qr.Size-finderPatternSize-4+j).SetReserved(bit)
			qr.Get(i, qr.Size-finderPatternSize-4+j).SetCell(bit)
			qr.Get(qr.Size-finderPatternSize-4+j, i).SetReserved(bit)
			qr.Get(qr.Size-finderPatternSize-4+j, i).SetCell(bit)
		}
	}
}

func (qr *QR) writeFormatInfo(mask int) {
	f := formatBits(mask)

	// Top-left
	for i := range 6 {
		qr.Get(finderPatternSize+1, i).SetCell((f>>i)&1 != 0)
	}

	qr.Get(finderPatternSize+1, finderPatternSize).SetCell((f>>6)&1 != 0)
	qr.Get(finderPatternSize+1, finderPatternSize+1).SetCell((f>>7)&1 != 0)
	qr.Get(finderPatternSize, finderPatternSize+1).SetCell((f>>8)&1 != 0)

	for i := 9; i < 15; i++ {
		qr.Get(14-i, finderPatternSize+1).SetCell((f>>i)&1 != 0)
	}

	// Top-right
	for i := 0; i < 8; i++ {
		qr.Get(qr.Size-1-i, finderPatternSize+1).SetCell((f>>i)&1 != 0)
	}

	// Bottom-left
	for i := 8; i < 15; i++ {
		qr.Get(finderPatternSize+1, qr.Size-finderPatternSize-8+i).SetCell((f>>i)&1 != 0)
	}
}
