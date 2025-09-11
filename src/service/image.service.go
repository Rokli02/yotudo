package service

import (
	"database/sql"
	"os"
	pathModule "path"
	"yotudo/src/database"
	"yotudo/src/database/entity"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
	"yotudo/src/settings"
)

type ImageService struct{}

var GlobalImageService *ImageService = nil

func (s *ImageService) Save(image *model.PossiblyNewImage) (*model.Image, error) {
	savedImage, err := repository.GlobalImageRepository.SaveOne(nil, image)
	if err != nil {
		return nil, err
	}

	return (&model.Image{}).FromEntity(savedImage), nil
}

func (s *ImageService) GetByName(name string) (*model.Image, error) {
	foundImage, err := repository.GlobalImageRepository.FindByName(nil, name)
	if err != nil {
		return nil, err
	}

	return (&model.Image{}).FromEntity(foundImage), nil
}

func (s *ImageService) GetBySource(source string) (*model.Image, error) {
	foundImage, err := repository.GlobalImageRepository.FindBySource(nil, source)
	if err != nil {
		return nil, err
	}

	return (&model.Image{}).FromEntity(foundImage), nil
}

func (s *ImageService) DeleteImageIfUnused(id int64) error {
	var foundImage *entity.Image

	if _foundImage, err := repository.GlobalImageRepository.FindById(nil, id); err != nil {
		return err
	} else {
		foundImage = _foundImage
	}

	tx, err := database.Instance.Begin()
	if err != nil {
		logger.Error("Couldn't start transaction for 'DecreaseReferedCountAndDelete'")
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			logger.Error(err)
		}
	}()

	if foundImage.ReferedCount <= 0 {
		logger.Debug("Have to delete image fron DB and FS")
		if err := repository.GlobalImageRepository.DeleteById(tx, id); err != nil {
			return err
		}

		logger.Debug("Deleted image fron DB")
		if err := os.Remove(pathModule.Join(settings.Global.App.ImagesLocation, foundImage.Name)); err != nil {
			logger.Error(err)
		}

		logger.Debug("Deleted image fron FS")
	}

	if err := tx.Commit(); err != nil {
		logger.Error("Couldn't commit started transaction for 'DecreaseReferedCountAndDelete'")
		return err
	}

	return nil
}
