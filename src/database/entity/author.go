package entity

type Author struct {
	Id   int64
	Name string
}

var _ Entity = Author{}

func (a Author) Migration(currentVersion MigrationVersion) []Migration {
	migrations := []Migration{}

	return MigrationsByVersion(migrations, currentVersion)
}

func (m Author) Template() string {
	return `
	CREATE TABLE author (
		id     INTEGER  PRIMARY KEY,
		name   TEXT     NOT NULL UNIQUE
	);

	CREATE INDEX author_name_index ON author(name);
	`
}
