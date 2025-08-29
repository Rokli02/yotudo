package entity

type Status struct {
	Id          int
	Name        string
	Description string
}

var _ Entity = Status{}

func (s Status) Migration(currentVersion MigrationVersion) []Migration {
	migrations := []Migration{}

	return MigrationsByVersion(migrations, currentVersion)
}

func (s Status) Template() string {
	return `
	CREATE TABLE status (
		id              INTEGER     PRIMARY KEY,
		name            TEXT        NOT NULL,
		description     TEXT
	);

	INSERT INTO status (id, name, description) VALUES(0, 'Letöltésre vár', 'Az alábbi videó már hozzá lett adva az adatbázishoz és letöltésre vár.');
	INSERT INTO status (name, description) VALUES('Folyamatban', 'Éppen folyamatban van a letöltés.'), ('Letöltve', 'A videó már letöltésre került');
	`
}
