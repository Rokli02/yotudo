package service

import (
	"yotudo/src/database/errors"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
)

type AuthorService struct{}

var GlobalAuthorService *AuthorService = nil

func (c *AuthorService) GetManyByPagination(filter string, page *model.Page, sort []model.Sort) *model.Pagination[[]model.Author] {
	authors, totalCount := repository.GlobalAuthorRepository.FindByPage(filter, page, sort)

	return &model.Pagination[[]model.Author]{
		Data:  authors,
		Count: totalCount,
	}
}

func (c *AuthorService) Save(newAuthorName string) (*model.Author, error) {
	return repository.GlobalAuthorRepository.SaveOne(newAuthorName)
}

func (c *AuthorService) Delete(id int64) (bool, error) {
	if repository.GlobalAuthorRepository.IsReferencingToMusic(id) {
		logger.Warning("Unable to delete Author, because it was used in a music, or contributor records")

		return false, errors.ErrUnableToDelete
	}

	return repository.GlobalAuthorRepository.DeleteOne(id), nil
}
