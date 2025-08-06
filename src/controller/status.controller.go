package controller

import (
	"yotudo/src/database/repository"
	"yotudo/src/model"
)

type StatusController struct {
	statusRepository *repository.Status
}

func NewStatusController(statusRepository *repository.Status) *StatusController {
	return &StatusController{statusRepository: statusRepository}
}

func (c *StatusController) GetAll() []model.Status {
	entities := c.statusRepository.FindAll()

	result := make([]model.Status, len(entities))
	for i, entity := range entities {
		result[i] = model.Status{
			Id:          entity.Id,
			Name:        entity.Name,
			Description: entity.Description,
		}
	}

	return result
}
