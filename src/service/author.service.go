package service

import (
	"yotudo/src/database/errors"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
)

type AuthorService struct {
	authorRepository *repository.Author
}

func NewAuthorService(authorRepository *repository.Author) *AuthorService {
	return &AuthorService{authorRepository: authorRepository}
}

func (c *AuthorService) GetManyByPagination(filter string, page *model.Page, sort []model.Sort) *model.Pagination[[]model.Author] {
	authors, totalCount := c.authorRepository.FindByPage(filter, page, sort)

	return &model.Pagination[[]model.Author]{
		Data:  authors,
		Count: totalCount,
	}
}

func (c *AuthorService) Save(newAuthorName string) (*model.Author, error) {
	return c.authorRepository.SaveOne(newAuthorName)
}

func (c *AuthorService) Delete(id int64) (bool, error) {
	if c.authorRepository.IsReferencingToMusic(id) {
		logger.Warning("Unable to delete Author, because it was used in a music, or contributor records")

		return false, errors.ErrUnableToDelete
	}

	return c.authorRepository.DeleteOne(id), nil
}
