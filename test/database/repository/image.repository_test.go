package repository_test

import (
	"testing"
	"yotudo/src/database/entity"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
)

func TestSaveOne(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	imageRepository := repository.NewImageRepository(db.Conn)

	savedImage, err := imageRepository.SaveOne(db.Conn, &model.PossiblyNewImage{Name: "kitalalt.jpg"})
	if err != nil {
		t.Error(err)
		return
	} else if savedImage.Id == 0 {
		t.Error("returned \"saved\" image didn't have a valid id.")
		return
	}

	row := db.Conn.QueryRow("SELECT id, path, referedCount FROM image WHERE id=?", savedImage.Id)
	image := entity.Image{}
	if err := row.Scan(&image.Id, &image.Name, &image.ReferedCount); err != nil {
		t.Error(err)
		return
	}

	logger.Debug(image.String())
}

func TestSaveOneWithNilConn(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	imageRepository := repository.NewImageRepository(db.Conn)

	savedImage, err := imageRepository.SaveOne(nil, &model.PossiblyNewImage{Name: "kitalalt.jpg"})
	if err != nil {
		t.Error(err)
		return
	} else if savedImage.Id == 0 {
		t.Error("returned \"saved\" image didn't have a valid id.")
		return
	}

	row := db.Conn.QueryRow("SELECT id, path, referedCount FROM image WHERE id=?", savedImage.Id)
	image := entity.Image{}
	if err := row.Scan(&image.Id, &image.Name, &image.ReferedCount); err != nil {
		t.Error(err)
		return
	}

	logger.Debug(image.String())
}

func TestFindById(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	imageRepository := repository.NewImageRepository(db.Conn)

	Must(imageRepository.SaveOne(nil, &model.PossiblyNewImage{Name: "test1.jpg"}))
	si2 := Must(imageRepository.SaveOne(nil, &model.PossiblyNewImage{Name: "karcsi.jpg"}))

	img, err := imageRepository.FindById(nil, si2.Id)
	if err != nil {
		t.Error(err)
		return
	}

	logger.Debug(img)
}

func TestFindByIdNotFound(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	imageRepository := repository.NewImageRepository(db.Conn)

	_, err := imageRepository.FindById(nil, 100)
	if err == nil {
		t.Error("Found image with id='100', even though it should not have")
		return
	}
}
