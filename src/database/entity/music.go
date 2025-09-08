package entity

import "encoding/json"

type Music struct {
	Id        int64
	Name      string
	Published *int
	Album     *string
	Url       string
	Filename  *string
	ImageId   *int64
	// Id of the current status of music
	Status int8
	// Id of the genre of music
	GenreId int
	// Id of the author of music
	AuthorId  int64
	UpdatedAt string
}

var _ Entity = Music{}

func (m Music) Migration(currentVersion MigrationVersion) []Migration {
	migrations := []Migration{
		{
			Version:         MigrationVersion{1, 1, 0},
			Migration:       "INSERT INTO image(name, referedCount) SELECT pic_filename as name, 1 as referedCount FROM music WHERE pic_filename IS NOT NULL;",
			SkipOnFreshInit: true,
		},
		{
			Version: MigrationVersion{1, 1, 0},
			Migration: `
			CREATE TABLE music_new (
				id              INTEGER     PRIMARY KEY,
				author_id       INTEGER     NOT NULL,
				name            TEXT        NOT NULL,
				published       INT,
				album           TEXT,
				genre_id        INTEGER     NOT NULL,
				url             TEXT        NOT NULL,
				filename        TEXT,
				pic_filename	TEXT,
				image_id		INTEGER,
				status          TINYINT     DEFAULT 0,
				updated_at      TEXT        NOT NULL,
				FOREIGN KEY(author_id)  REFERENCES author(id),
				FOREIGN KEY(genre_id)   REFERENCES genre(id),
				FOREIGN KEY(image_id)   REFERENCES image(id)
			);

			INSERT INTO music_new(id, author_id, name, published, album, genre_id, url, filename, pic_filename, status, updated_at) SELECT id, author_id, name, published, album, genre_id, url, filename, pic_filename, status, updated_at FROM music;

			DROP TABLE music;

			ALTER TABLE music_new RENAME TO music;

			UPDATE music SET image_id=image.id FROM image WHERE music.pic_filename=image.name;

			ALTER TABLE music DROP COLUMN pic_filename;`,
			SkipOnFreshInit: true,
		},
		{
			Version: MigrationVersion{1, 1, 1},
			Migration: `
			CREATE TRIGGER IF NOT EXISTS increase_image_ref_on_insert AFTER INSERT ON music FOR EACH ROW
			WHEN new.image_id IS NOT NULL
			BEGIN
				UPDATE image SET referedCount = referedCount + 1 WHERE id = new.image_id;
			END;

			CREATE TRIGGER IF NOT EXISTS decrease_image_ref_on_delete AFTER DELETE ON music FOR EACH ROW
			WHEN old.image_id IS NOT NULL
			BEGIN
				UPDATE image SET referedCount = referedCount - 1 WHERE id = old.image_id;
			END;

			CREATE TRIGGER IF NOT EXISTS vary_image_ref_on_update AFTER UPDATE OF image_id ON music FOR EACH ROW
			WHEN (old.image_id IS NOT NULL OR new.image_id IS NOT NULL) AND old.image_id IS NOT new.image_id
			BEGIN
				UPDATE image SET referedCount = referedCount - 1 WHERE id = old.image_id;
				UPDATE image SET referedCount = referedCount + 1 WHERE id = new.image_id;
			END;
			`,
		},
	}

	return MigrationsByVersion(migrations, currentVersion)
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
		image_id		INTEGER,
		status          TINYINT     DEFAULT 0,
		updated_at      TEXT        NOT NULL,
		FOREIGN KEY(author_id)  REFERENCES author(id),
		FOREIGN KEY(genre_id)   REFERENCES genre(id),
		FOREIGN KEY(image_id)   REFERENCES image(id)
	);`
}

func (m *Music) String() string {
	jb, _ := json.Marshal(m)
	return string(jb)
}

func (m *Music) FromScan(Scan func(dest ...any) error) (*Music, error) {
	return m, Scan(&m.Id, &m.AuthorId, &m.Name, &m.Published, &m.Album, &m.GenreId, &m.Url, &m.Filename, &m.ImageId, &m.Status, &m.UpdatedAt)
}
