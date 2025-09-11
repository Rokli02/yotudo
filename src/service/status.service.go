package service

import (
	"yotudo/src/database/repository"
	"yotudo/src/model"
)

type StatusService struct{}

var GlobalStatusService *StatusService = nil

func (c *StatusService) GetAll() []model.Status {
	entities := repository.GlobalStatusRepository.FindAll()

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
