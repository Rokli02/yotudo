package file

import (
	"strings"
	"testing"
	"yotudo/src/lib/logger"
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
	service := file.NewFileService()

	purgedFileName := service.PurgeFileName(fileName)

	logger.Debug("Purged filename", purgedFileName)

	if len(purgedFileName) == len(fileName) || purgedFileName == fileName || strings.ContainsAny(purgedFileName, "|:") {
		t.Error("Filename purge didn't get executed")
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
