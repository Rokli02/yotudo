package qrcode

import (
	"errors"
	"strings"
	"unicode"
)

type metadata struct {
	Version   version
	Mode      encodingMode
	CharCount uint32
}

func calculateMetadata(input string) (*metadata, error) {
	metadata := new(metadata)

	metadata.Mode = determineMode(input)
	metadata.CharCount = countCharacters(input, metadata.Mode)

	if version, err := metadata.determineVersion(); err != nil {
		return nil, err
	} else {
		metadata.Version = version
	}

	return metadata, nil
}

type version uint32

type encodingMode int

func (m encodingMode) indicator() byte {
	switch m {
	case ModeNumeric:
		return 0b0001
	case ModeAlphanumeric:
		return 0b0010
	default:
		return 0b0100
	}
}

const (
	ModeNumeric encodingMode = iota
	ModeAlphanumeric
	ModeByte
)

var capacityTableByModes = map[encodingMode][]uint32{
	ModeNumeric: {
		41, 77, 127, 187, 255, 322, 370, 461, 552, 652,
		772, 883, 1022, 1101, 1250, 1408, 1548, 1725, 1903, 2061,
		2232, 2409, 2620, 2812, 3057, 3283, 3517, 3669, 3909, 4158,
		4417, 4686, 4965, 5253, 5529, 5836, 6153, 6479, 6743, 7089,
	},
	ModeAlphanumeric: {
		25, 47, 77, 114, 154, 195, 224, 279, 335, 395,
		468, 535, 619, 667, 758, 854, 938, 1046, 1153, 1249,
		1352, 1460, 1588, 1704, 1853, 1990, 2132, 2223, 2369, 2520,
		2677, 2840, 3009, 3183, 3351, 3537, 3729, 3927, 4087, 4296,
	},
	ModeByte: {
		17, 32, 53, 78, 106, 134, 154, 192, 230, 271,
		321, 367, 425, 458, 520, 586, 644, 718, 792, 858,
		929, 1003, 1091, 1171, 1273, 1367, 1465, 1528, 1628, 1732,
		1840, 1952, 2068, 2188, 2303, 2431, 2563, 2699, 2809, 2953,
	},
}

func (m encodingMode) String() string {
	switch m {
	case ModeNumeric:
		return "Numeric"
	case ModeAlphanumeric:
		return "Alphanumeric"
	default:
		return "Byte"
	}
}

func countCharacters(input string, mode encodingMode) uint32 {
	if mode == ModeByte {
		return uint32(len([]byte(input)))
	}
	return uint32(len(input))
}

func determineMode(input string) encodingMode {
	isNumeric := true

	for _, r := range input {
		if unicode.IsDigit(r) {
			continue
		} else {
			isNumeric = false
		}

		if strings.ContainsRune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:", r) {
			continue
		} else {
			return ModeByte
		}
	}

	if isNumeric {
		return ModeNumeric
	}

	return ModeAlphanumeric
}

func (m *metadata) determineVersion() (version, error) {
	capacityTable := capacityTableByModes[m.Mode]

	for i, capacity := range capacityTable {
		if m.CharCount <= capacity {
			return version(i + 1), nil
		}
	}

	return 0, errors.New("input too large for QR version 40")
}

func (m *metadata) size() uint32 {
	return uint32(m.Version)*4 + 17
}
