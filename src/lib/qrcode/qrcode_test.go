package qrcode_test

import (
	"fmt"
	"testing"
	"yotudo/src/lib/qrcode"
)

func TestEncodeTextToQRCode(t *testing.T) {
	var text string = "MECSODA PRANK"
	qr, err := qrcode.Encode(text)
	if err != nil {
		t.Error(err)
		return
	}

	qrCode := qr.JsonBitmap()

	fmt.Println(qrCode)
}
