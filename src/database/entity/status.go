package entity

type Status struct {
	Id          int
	Name        string
	Description string
}

var _ Entity = Status{}

func (s Status) Migration(currentVersion MigrationVersion) []Migration {
	migrations := []Migration{
		{
			Version: MigrationVersion{0, 1, 0},
			Migration: `
			INSERT INTO status (id, name, description) VALUES(0, 'Letöltésre vár', 'Az alábbi videó már hozzá lett adva az adatbázishoz és letöltésre vár.');
			INSERT INTO status (name, description) VALUES('Folyamatban', 'Éppen folyamatban van a letöltés.'), ('Letöltve', 'A videó már letöltésre került');
			`,
		},
		{
			Version: MigrationVersion{1, 1, 1},
			Migration: `
				UPDATE status SET name='Feldolgozás', description='A hozzáadott zene készen áll a feldolgozásra' WHERE id = 0;
				UPDATE status SET name='Folyamatban', description='Az elindított folyamat - előfeldolgozása vagy letöltése - hamarosan elkészül' WHERE id = 1;
				UPDATE status SET name='Letöltés', description='A zene készen áll a letöltésre' WHERE id = 2;
			`,
		},
	}

	return MigrationsByVersion(migrations, currentVersion)
}

func (s Status) Template() string {
	return `
	CREATE TABLE status (
		id              INTEGER     PRIMARY KEY,
		name            TEXT        NOT NULL,
		description     TEXT
	);`
}
