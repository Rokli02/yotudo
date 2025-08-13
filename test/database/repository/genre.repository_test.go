package repository_test

import (
	"testing"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
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

func TestGenreIsAlreadyUsed(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	genreRepository := repository.NewGenreRepository(db.Conn)
	authorRepository := repository.NewAuthorRepository(db.Conn)
	musicRepository := repository.NewMusicRepository(db.Conn, repository.NewContributorRepository(db.Conn))

	var expectedGenreId int64

	genres := genreRepository.FindAll()
	expectedGenreId = genres[0].Id
	johnLenon, _ := authorRepository.SaveOne("John Lenon")
	musicRepository.SaveOne(&model.NewMusic{
		Name:    "Test Muzsik",
		Url:     "http://youtoube.com/watch?v=abc1234567g",
		Author:  model.OptionalAuthor{Id: &johnLenon.Id},
		GenreId: expectedGenreId,
	})

	if !genreRepository.IsAlreadyUsed(expectedGenreId) {
		t.Error("Genre is used in a music entity but did not found")

		return
	}
}

func TestGenreIsNotAlreadyUsed(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	genreRepository := repository.NewGenreRepository(db.Conn)
	authorRepository := repository.NewAuthorRepository(db.Conn)
	musicRepository := repository.NewMusicRepository(db.Conn, repository.NewContributorRepository(db.Conn))

	var expectedGenreId int64

	genres := genreRepository.FindAll()
	expectedGenreId = genres[1].Id
	johnLenon, _ := authorRepository.SaveOne("John Lenon")
	musicRepository.SaveOne(&model.NewMusic{
		Name:    "Test Muzsik",
		Url:     "http://youtoube.com/watch?v=abc1234567g",
		Author:  model.OptionalAuthor{Id: &johnLenon.Id},
		GenreId: genres[0].Id,
	})

	if genreRepository.IsAlreadyUsed(expectedGenreId) {
		t.Error("Genre is not used in a music entity but was found")

		return
	}
}
