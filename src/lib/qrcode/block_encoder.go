package qrcode

import (
	"yotudo/src/lib/qrcode/bitbuffer"
	"yotudo/src/lib/qrcode/rs"
)

func encodeBitBufferIntoBlocks(buffer *bitbuffer.BitBuffer, metadata *metadata) (*bitbuffer.BitBuffer, error) {
	blockCount := int(ecTableL[metadata.Version-1][0])
	ecBytes := int(ecTableL[metadata.Version-1][1])
	dataPerBlock := buffer.ByteLen() / uint(blockCount)
	blocks := make([]*bitbuffer.BitBuffer, 0, blockCount)

	for i, offset := 0, uint(0); i < blockCount; i++ {
		if slicedData, err := buffer.SliceBytes(offset, offset+dataPerBlock); err != nil {
			return nil, err
		} else {
			block := rs.Encode(slicedData, ecBytes)
			blocks = append(blocks, block)
			offset += dataPerBlock
		}
	}

	// interleave the blocks
	result := bitbuffer.New()

	// Kombájn détä bloaksz
	for i := uint(0); i < dataPerBlock*8; i += 8 {
		for _, block := range blocks {
			chunk, length, err := block.GetBitChunk(i, 8)
			if err != nil {
				return nil, err
			}

			if length > 0 {
				err = result.AppendByte(chunk[0], length)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	// Kombájn errór korreksön bloaksz
	for i, mégegykör := dataPerBlock*8, true; mégegykör; i += 8 {
		mégegykör = false

		for _, block := range blocks {

			if i >= block.Len() {
				continue
			}

			chunk, length, err := block.GetBitChunk(i, 8)
			if err != nil {
				return nil, err
			}

			if length > 0 {
				err = result.AppendByte(chunk[0], length)
				if err != nil {
					return nil, err
				}
			}

			mégegykör = true
		}
	}

	remainder := uint(remainderTable[metadata.Version])
	if remainder != 0 {
		result.AppendNBits(remainder, false)
	}

	return result, nil
}

var ecTableL = [][2]int8{
	{1, 7}, {1, 10}, {1, 15}, {1, 20}, {1, 26},
	{2, 18}, {2, 20}, {2, 24}, {2, 30}, {4, 18},
	{4, 20}, {4, 24}, {4, 26}, {4, 30}, {6, 22},
	{6, 24}, {6, 28}, {6, 30}, {7, 28}, {8, 28},
	{8, 30}, {9, 30}, {9, 30}, {10, 30}, {12, 26},
	{12, 28}, {12, 30}, {13, 30}, {14, 30}, {15, 30},
	{16, 30}, {17, 30}, {18, 30}, {19, 30}, {19, 30},
	{20, 30}, {21, 30}, {22, 30}, {24, 30}, {25, 30},
}

var remainderTable = [40]uint8{
	/*  1    */ 0,
	/*  2-6  */ 7, 7, 7, 7, 7,
	/*  7-13 */ 0, 0, 0, 0, 0, 0, 0,
	/* 14-20 */ 3, 3, 3, 3, 3, 3, 3,
	/* 21-27 */ 4, 4, 4, 4, 4, 4, 4,
	/* 28-34 */ 3, 3, 3, 3, 3, 3, 3,
	/* 35-40 */ 0, 0, 0, 0, 0, 0,
}
