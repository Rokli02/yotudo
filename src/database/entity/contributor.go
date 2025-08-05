package entity

type Contributor struct {
	MusicId  int64
	AuthorId int64
}

var _ Entity = Contributor{}

func (c Contributor) Migration(currentVersion [3]int) []Migration {
	migrations := []Migration{
		{
			Version:   [3]int{0, 1, 0},
			Migration: c.Template(),
		},
	}

	return migrationsOfVersion(migrations, currentVersion)
}

func (c Contributor) Template() string {
	return `
	CREATE TABLE contributor (
		music_id    INTEGER     NOT NULL,
		author_id   INTEGER     NOT NULL,
		UNIQUE(music_id, author_id) ON CONFLICT IGNORE,
		FOREIGN KEY(music_id)   REFERENCES music(id) ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY(author_id)   REFERENCES author(id) ON DELETE CASCADE ON UPDATE CASCADE
	);
	`
}
