package service

import (
	"yotudo/src/database/repository"
	"yotudo/src/model"
)

type StatusService struct {
	statusRepository *repository.Status
}

func NewStatusService(statusRepository *repository.Status) *StatusService {
	return &StatusService{statusRepository: statusRepository}
}

func (c *StatusService) GetAll() []model.Status {
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
