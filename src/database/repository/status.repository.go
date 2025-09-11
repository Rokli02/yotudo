package repository

import (
	"yotudo/src/database"
	"yotudo/src/database/entity"
	"yotudo/src/lib/logger"
)

type StatusRepository struct{}

var GlobalStatusRepository *StatusRepository = nil

func (s *StatusRepository) FindAll() []entity.Status {
	var statuses []entity.Status

	rows, err := database.Instance.Query("SELECT id, name, description FROM status;")
	if err != nil {
		logger.Error(err)

		return []entity.Status{}
	}

	statuses = make([]entity.Status, 0, 3)

	defer rows.Close()
	for rows.Next() {
		status := entity.Status{}

		err := rows.Scan(&status.Id, &status.Name, &status.Description)
		if err != nil {
			logger.Warning("Status.FindAll:", err)
		} else {
			statuses = append(statuses, status)
		}
	}

	return statuses
}
