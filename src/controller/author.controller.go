package controller

import (
	"yotudo/src/database/errors"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
)

type AuthorController struct {
	authorRepository *repository.Author
}

func NewAuthorController(authorRepository *repository.Author) *AuthorController {
	return &AuthorController{authorRepository: authorRepository}
}

func (c *AuthorController) GetManyByPagination(filter string, page *model.Page, sort []model.Sort) *model.Pagination[[]model.Author] {
	authors, totalCount := c.authorRepository.FindByPage(filter, page, sort)

	return &model.Pagination[[]model.Author]{
		Data:  authors,
		Count: totalCount,
	}
}

func (c *AuthorController) Save(newAuthorName string) (*model.Author, error) {
	return c.authorRepository.SaveOne(newAuthorName)
}

func (c *AuthorController) Delete(id int64) (bool, error) {
	if c.authorRepository.IsReferencingToMusic(id) {
		logger.Warning("Unable to delete Author, because it was used in a music, or contributor records")

		return false, errors.ErrUnableToDelete
	}

	return c.authorRepository.DeleteOne(id), nil
}
