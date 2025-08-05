package entity

type Genre struct {
	Id   int64
	Name string
}

var _ Entity = Genre{}

func (g Genre) Migration(currentVersion [3]int) []Migration {
	migrations := []Migration{
		{
			Version:   [3]int{0, 1, 0},
			Migration: g.Template(),
		},
	}

	return migrationsOfVersion(migrations, currentVersion)
}

func (g Genre) Template() string {
	return `
	CREATE TABLE genre (
		id      INTEGER     PRIMARY KEY,
		name    TEXT        NOT NULL UNIQUE
	);

	INSERT INTO genre (name) VALUES('Ismeretlen'), ('Rock'), ('Pop'), ('Rap'), ('Metál');
	`
}
