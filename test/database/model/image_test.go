package model_test

import (
	"testing"
	"yotudo/src/model"
)

func TestPossiblyNewImageHasPath(t *testing.T) {
	var img *model.PossiblyNewImage = &model.PossiblyNewImage{Name: "valami.jpg"}

	if !img.HasName() {
		t.Error("Image should have had a path")
	}
}

func TestPossiblyNewImageDoesNotHavePath(t *testing.T) {
	var img *model.PossiblyNewImage = &model.PossiblyNewImage{}

	if img.HasName() {
		t.Error("Image should not have had a path")
	}
}

func TestPossiblyNewImageIsNilHasPath(t *testing.T) {
	var img *model.PossiblyNewImage = nil

	if img.HasName() {
		t.Error("Image should not have had a path")
	}
}
