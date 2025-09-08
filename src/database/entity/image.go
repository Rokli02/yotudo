package entity

import "fmt"

type Image struct {
	Id           int64
	Name         string
	Source       *string
	ReferedCount int
}

var _ Entity = Image{}

func (i *Image) String() string {
	return fmt.Sprintf("Image(id=%d, path=%s, referedCount=%d)", i.Id, i.Name, i.ReferedCount)
}

func (i Image) Template() string {
	return `
	CREATE TABLE IF NOT EXISTS image (
		id INTEGER PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		sourcePath TEXT,
		referedCount INTEGER DEFAULT 0
	);
	`
}

func (i Image) Migration(currentVersion MigrationVersion) []Migration {
	migrations := []Migration{
		{
			Version:   MigrationVersion{1, 1, 0},
			Migration: i.Template(),
		},
		{
			Version:   MigrationVersion{1, 1, 0},
			Migration: "INSERT INTO image(name, referedCount) SELECT pic_filename as name, 1 as referedCount FROM music WHERE pic_filename IS NOT NULL;",
		},
	}

	return MigrationsByVersion(migrations, currentVersion)
}
