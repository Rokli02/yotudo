package qrcode

import (
	"fmt"
	"strings"
)

type QR struct {
	Size int
	Data []Cell
}

func Encode(text string) (*QR, error) {
	metadata, err := calculateMetadata(text)
	if err != nil {
		return nil, err
	}

	buffer, err := encodeDataIntoBitBuffer([]byte(text), metadata)
	if err != nil {
		return nil, err
	}

	buffer, err = encodeBitBufferIntoBlocks(buffer, metadata)
	if err != nil {
		return nil, err
	}

	qr, _ := createQRWithBestMask(buffer, metadata)

	return qr, nil
}

func (qr *QR) AsciiString() string {
	sb := strings.Builder{}

	y1 := 0
	for i, cell := range qr.Data {
		y2 := i / qr.Size

		if y1 != y2 {
			sb.WriteByte('\n')
			y1 = y2
		}

		if cell.IsSet() {
			sb.WriteRune('█')
		} else {
			sb.WriteRune(' ')
		}
	}

	return sb.String()
}

func (qr *QR) Json() string {
	sb := strings.Builder{}

	for i, data := range qr.Data {
		if data.IsSet() {
			sb.WriteString("true")
		} else {
			sb.WriteString("false")
		}
		if i != len(qr.Data)-1 {
			sb.WriteString(", ")
		}
	}

	return fmt.Sprintf("{ size: %d, data: [%s]}", qr.Size, sb.String())
}

func (qr *QR) JsonBitmap() string {
	sb := strings.Builder{}

	for i, data := range qr.Data {
		sb.WriteString(fmt.Sprintf("%d", data))

		if i != len(qr.Data)-1 {
			sb.WriteString(", ")
		}
	}

	return fmt.Sprintf("{ size: %d, attribute: %s, data: [%s]}", qr.Size, qr.AttributesJson(), sb.String())
}
