package service

import (
	"yotudo/src/database/errors"
	"yotudo/src/database/repository"
	"yotudo/src/model"
)

type GenreService struct {
	genreRepository *repository.Genre
}

func NewGenreService(genreRepository *repository.Genre) *GenreService {
	return &GenreService{genreRepository: genreRepository}
}

func (c *GenreService) GetAll() []model.Genre {
	entities := c.genreRepository.FindAll()

	genres := make([]model.Genre, len(entities))
	for i, entity := range entities {
		genres[i] = model.Genre{Id: entity.Id, Name: entity.Name}
	}

	return genres
}

func (c *GenreService) Save(genreName string) (*model.Genre, error) {
	entity, err := c.genreRepository.SaveOne(genreName)
	if err != nil {
		return nil, err
	}

	return &model.Genre{Id: entity.Id, Name: entity.Name}, nil
}

func (c *GenreService) Rename(id int64, newGenreName string) (*model.Genre, error) {
	entity, err := c.genreRepository.Rename(id, newGenreName)
	if err != nil {
		return nil, err
	}

	return &model.Genre{Id: entity.Id, Name: entity.Name}, nil
}

func (c *GenreService) Delete(id int64) error {
	if c.genreRepository.IsAlreadyUsed(id) {
		return errors.ErrUnableToDelete
	}

	return c.genreRepository.DeleteOne(id)
}
