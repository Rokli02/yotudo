package repository_test

import (
	"fmt"
	"runtime"
	"yotudo/src/database"
)

func getInMemoryDB(shared ...bool) *database.Database {
	return database.LoadDatabase(func(opts *database.DatabaseOptions) {
		if len(shared) > 0 && shared[0] {
			opts.SetLocation("file::memory:?cache=shared")
		} else {
			opts.SetLocation(":memory:?cache=shared")
		}
	})
}

func Must[T any](value T, err error) T {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		panic(fmt.Sprintf("Error in \"%s:%d\": %s\n", file, line, err.Error()))
	}

	return value
}
