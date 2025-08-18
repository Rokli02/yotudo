package service

import (
	"strings"
	"testing"
	"yotudo/src/model"
	"yotudo/src/service"
)

func TestValidName(t *testing.T) {
	fileName := "valami jo.exe"
	service := service.NewFileService()

	if !service.ValidName(fileName) {
		t.Errorf("\"%s\"is a valid filename", fileName)
	}
}

func TestBlankStringAsName(t *testing.T) {
	fileName := "       "
	service := service.NewFileService()

	if service.ValidName(fileName) {
		t.Errorf("\"%s\"is not a valid filename", fileName)
	}
}

func TestInvalidFilename(t *testing.T) {
	fileName := "nem|jo:.exe"
	service := service.NewFileService()

	if service.ValidName(fileName) {
		t.Errorf("\"%s\"is not a valid filename", fileName)
	}
}

func TestPurgeFilename(t *testing.T) {
	fileName := " nem|jo:.exe"
	expectedName := "nem_jo_.exe"
	service := service.NewFileService()

	purgedFileName := service.PurgeFileName(fileName)

	if strings.ContainsAny(purgedFileName, "|:") || purgedFileName != expectedName {
		t.Errorf("Filename purge didn't get executed (expected \"%s\", but got \"%s\")", expectedName, purgedFileName)
	}
}

func TestPurgeFilenameNotNeeded(t *testing.T) {
	fileName := "jo fajlnev.png"
	service := service.NewFileService()

	purgedFileName := service.PurgeFileName(fileName)

	if purgedFileName != fileName {
		t.Error("Unnecessary filename purge was executed")
	}
}

func TestCreateFilename(t *testing.T) {
	music := &model.Music{
		Name:   "test",
		Author: model.Author{Name: "test.author"},
		Contributors: []model.Author{
			{Name: "sanesz"},
			{Name: "kanesz"},
		},
	}
	expectedFilenamePrefix := "test.author (sanesz, kanesz) - test ["
	service := service.NewFileService()

	filename := service.CreateFilename(music)

	if !strings.HasPrefix(filename, expectedFilenamePrefix) {
		t.Errorf("Unexpected filename was created (%s)", filename)

		return
	}
}

func TestCreateFilenameWithPurge(t *testing.T) {
	music := &model.Music{
		Name:   "test\\/",
		Author: model.Author{Name: "test.author:"},
		Contributors: []model.Author{
			{Name: "sanesz"},
			{Name: "kanesz"},
		},
	}
	expectedFilenamePrefix := "test.author_ (sanesz, kanesz) - test__ ["
	service := service.NewFileService()

	filename := service.CreateFilename(music)

	if !strings.HasPrefix(filename, expectedFilenamePrefix) {
		t.Errorf("Unexpected filename was created (%s)", filename)

		return
	}
}
