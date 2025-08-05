package repository

import (
	"database/sql"
	"yotudo/src/database/entity"
	"yotudo/src/lib/logger"
)

type Status struct {
	db *sql.DB
}

func NewStatusRepository(db *sql.DB) *Status {
	return &Status{db: db}
}

func (s *Status) FindAll() []entity.Status {
	var statuses []entity.Status

	rows, err := s.db.Query("SELECT id, name, description FROM status;")
	if err != nil {
		logger.Warning(err)

		return []entity.Status{}
	}

	statuses = make([]entity.Status, 0, 3)

	defer rows.Close()
	for rows.Next() {
		status := entity.Status{}

		err := rows.Scan(&status.Id, &status.Name, &status.Description)
		if err != nil {
			logger.Warning(err)
		} else {
			statuses = append(statuses, status)
		}
	}

	return statuses
}
