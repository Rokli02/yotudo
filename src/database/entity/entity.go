package entity

type Entity interface {
	Template() string
	Migration(currentVersion [3]int) []Migration
}

type Migration struct {
	Version   [3]int
	Migration string
}

func migrationsOfVersion(migrations []Migration, currentVersion [3]int) []Migration {
	for i := len(migrations) - 1; i >= 0; i-- {
		vers := migrations[i].Version

		if vers[0] <= currentVersion[0] && vers[1] <= currentVersion[1] && vers[2] <= currentVersion[2] {
			return migrations[i+1:]
		}
	}

	return nil
}
