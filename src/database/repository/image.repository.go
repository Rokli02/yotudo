package repository

import (
	"yotudo/src/database"
	"yotudo/src/database/entity"
	"yotudo/src/database/errors"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
)

type ImageRepository struct{}

var GlobalImageRepository *ImageRepository = nil

func (i *ImageRepository) FindById(conn Connection, id int64) (*entity.Image, error) {
	if conn == nil {
		conn = database.Instance
	}

	img := &entity.Image{}
	if err := conn.QueryRow("SELECT id, name, sourcePath, referedCount FROM image WHERE id=? LIMIT 1;", id).Scan(&img.Id, &img.Name, &img.Source, &img.ReferedCount); err != nil {
		logger.Warning("Image.FindById:", err)

		return nil, errors.ErrNotFound
	}

	return img, nil
}

func (i *ImageRepository) FindByName(conn Connection, name string) (*entity.Image, error) {
	if conn == nil {
		conn = database.Instance
	}

	img := &entity.Image{}
	if err := conn.QueryRow("SELECT id, name, sourcePath, referedCount FROM image WHERE name=? LIMIT 1;", name).Scan(&img.Id, &img.Name, &img.Source, &img.ReferedCount); err != nil {
		logger.Warning("Image.FindById:", err)

		return nil, errors.ErrNotFound
	}

	return img, nil
}

func (i *ImageRepository) FindBySource(conn Connection, source string) (*entity.Image, error) {
	if conn == nil {
		conn = database.Instance
	}

	img := &entity.Image{}
	if err := conn.QueryRow("SELECT id, name, sourcePath, referedCount FROM image WHERE sourcePath=? LIMIT 1;", source).Scan(&img.Id, &img.Name, &img.Source, &img.ReferedCount); err != nil {
		logger.Warning("Image.FindById:", err)

		return nil, errors.ErrNotFound
	}

	return img, nil
}

func (i *ImageRepository) SaveOne(conn Connection, image *model.PossiblyNewImage) (*entity.Image, error) {
	if conn == nil {
		conn = database.Instance
	}

	if res, err := conn.Exec("INSERT INTO image (name, sourcePath) VALUES (?, ?);", image.Name, image.Source); err != nil {
		logger.Error(err)
		return nil, errors.ErrUnableToSave
	} else if id, err := res.LastInsertId(); err != nil {
		logger.Error(err)
		return nil, errors.ErrUnableToSave
	} else {
		return &entity.Image{Id: id, Name: image.Name, ReferedCount: 1}, nil
	}
}

func (i *ImageRepository) DeleteById(conn Connection, id int64) error {
	if conn == nil {
		conn = database.Instance
	}

	if _, err := conn.Exec("DELETE FROM image WHERE id=?;", id); err != nil {
		logger.Error(err)
		return errors.ErrUnableToDelete
	}

	return nil
}
