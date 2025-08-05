package repository_test

import (
	"testing"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
)

func TestGetAllStatuses(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	repo := repository.NewStatusRepository(db.Conn)

	allStatus := repo.FindAll()
	logger.Info("Statuses:", allStatus)
}
