package entity

type Music struct {
	Id          int64
	Name        string
	Published   *int
	Album       *string
	Url         string
	Filename    *string
	PicFilename *string
	// Id of the current status of music
	Status int8
	// Id of the genre of music
	GenreId int
	// Id of the author of music
	AuthorId  int64
	UpdatedAt string
}

var _ Entity = Music{}

func (m Music) Migration(currentVersion [3]int) []Migration {
	migrations := []Migration{
		{
			Version:   [3]int{0, 1, 0},
			Migration: m.Template(),
		},
	}

	return migrationsOfVersion(migrations, currentVersion)
}

func (m Music) Template() string {
	return `CREATE TABLE music (
		id              INTEGER     PRIMARY KEY,
		author_id       INTEGER     NOT NULL,
		name            TEXT        NOT NULL,
		published       INT,
		album           TEXT,
		genre_id        INTEGER     NOT NULL,
		url             TEXT        NOT NULL,
		filename        TEXT,
		pic_filename	TEXT,
		status          TINYINT     DEFAULT 0,
		updated_at      TEXT        NOT NULL,
		FOREIGN KEY(author_id)  REFERENCES author(id),
		FOREIGN KEY(genre_id)   REFERENCES genre(id)
	);`
}
