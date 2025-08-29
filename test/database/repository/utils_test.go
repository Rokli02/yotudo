package repository_test

import "yotudo/src/database"

func getInMemoryDB(shared ...bool) *database.Database {
	return database.LoadDatabase(func(opts *database.DatabaseOptions) {
		if len(shared) > 0 && shared[0] {
			opts.SetLocation("file::memory:?cache=shared")
		} else {
			opts.SetLocation(":memory:?cache=shared")
		}
	})
}
