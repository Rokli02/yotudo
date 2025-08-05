package repository_test

import "yotudo/src/database"

func getInMemoryDB() *database.Database {
	return database.LoadDatabase(func(opts *database.DatabaseOptions) {
		opts.SetLocation(":memory:")
	})
}
