package database

import (
	"database/sql"
	"yotudo/src/database/repository"
)

var InfoRepository *repository.Info
var StatusRepository *repository.Status
var GenreRepository *repository.Genre
var AuthorRepository *repository.Author
var ContributorRepository *repository.Contributor
var MusicRepository *repository.Music

func InitRepositories(conn *sql.DB) {
	InfoRepository = repository.NewInfoRepository(conn)
	StatusRepository = repository.NewStatusRepository(conn)
	GenreRepository = repository.NewGenreRepository(conn)
	AuthorRepository = repository.NewAuthorRepository(conn)
	ContributorRepository = repository.NewContributorRepository(conn)
	MusicRepository = repository.NewMusicRepository(conn)
}
