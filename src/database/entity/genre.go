package entity

type Genre struct {
	Id   int64
	Name string
}

var _ Entity = Genre{}

func (g Genre) Migration(currentVersion MigrationVersion) []Migration {
	migrations := []Migration{}

	return MigrationsByVersion(migrations, currentVersion)
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
