package file

import (
	"strings"
	"testing"
	"yotudo/src/service/file"
)

func TestValidName(t *testing.T) {
	fileName := "valami jo.exe"
	service := file.NewFileService()

	if !service.ValidName(fileName) {
		t.Errorf("\"%s\"is a valid filename", fileName)
	}
}

func TestBlankStringAsName(t *testing.T) {
	fileName := "       "
	service := file.NewFileService()

	if service.ValidName(fileName) {
		t.Errorf("\"%s\"is not a valid filename", fileName)
	}
}

func TestInvalidFilename(t *testing.T) {
	fileName := "nem|jo:.exe"
	service := file.NewFileService()

	if service.ValidName(fileName) {
		t.Errorf("\"%s\"is not a valid filename", fileName)
	}
}

func TestPurgeFilename(t *testing.T) {
	fileName := " nem|jo:.exe"
	expectedName := "nem_jo_.exe"
	service := file.NewFileService()

	purgedFileName := service.PurgeFileName(fileName)

	if strings.ContainsAny(purgedFileName, "|:") || purgedFileName != expectedName {
		t.Errorf("Filename purge didn't get executed (expected \"%s\", but got \"%s\")", expectedName, purgedFileName)
	}
}

func TestPurgeFilenameNotNeeded(t *testing.T) {
	fileName := "jo fajlnev.png"
	service := file.NewFileService()

	purgedFileName := service.PurgeFileName(fileName)

	if purgedFileName != fileName {
		t.Error("Unnecessary filename purge was executed")
	}
}
