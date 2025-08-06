package controller

import (
	"yotudo/src/database/repository"
	"yotudo/src/model"
)

type GenreController struct {
	genreRepository *repository.Genre
}

func NewGenreController(genreRepository *repository.Genre) *GenreController {
	return &GenreController{genreRepository: genreRepository}
}

func (c *GenreController) GetAll() []model.Genre {
	entities := c.genreRepository.FindAll()

	genres := make([]model.Genre, len(entities))
	for i, entity := range entities {
		genres[i] = model.Genre{Id: entity.Id, Name: entity.Name}
	}

	return genres
}

func (c *GenreController) Save(genreName string) (*model.Genre, error) {
	entity, err := c.genreRepository.SaveOne(genreName)
	if err != nil {
		return nil, err
	}

	return &model.Genre{Id: entity.Id, Name: entity.Name}, nil
}

func (c *GenreController) Rename(id int64, newGenreName string) (*model.Genre, error) {
	entity, err := c.genreRepository.Rename(id, newGenreName)
	if err != nil {
		return nil, err
	}

	return &model.Genre{Id: entity.Id, Name: entity.Name}, nil
}
