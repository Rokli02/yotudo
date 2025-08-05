package repository_test

import (
	"testing"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
)

func TestGenreFindAll(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	repo := repository.NewGenreRepository(db.Conn)

	allGenre := repo.FindAll()
	logger.Info("All Genre:", allGenre)
}

func TestGenreSave(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	repo := repository.NewGenreRepository(db.Conn)

	savedGenre, err := repo.SaveOne("TestKat")
	logger.Info("Saved Genre:", savedGenre)

	if err != nil {
		t.Error(err)
	}
}

func TestGenreRename(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	repo := repository.NewGenreRepository(db.Conn)

	allGenre := repo.FindAll()
	logger.Info("All Genre:", allGenre)

	_, err := repo.Rename(1, "Közismert")

	if err != nil {
		t.Error(err)
	}

	allGenre = repo.FindAll()
	logger.Info("All Genre After Rename:", allGenre)
}
