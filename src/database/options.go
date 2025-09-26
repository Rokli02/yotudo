package database

import (
	"os"
	"path"
	settingsModule "yotudo/src/settings"
)

type DatabaseOptions struct {
	location string
}

type DatabaseOptionsFunc func(opts *DatabaseOptions)

func DefaultDatabaseOptions(settings settingsModule.DatabaseSettings) *DatabaseOptions {
	wd, _ := os.Getwd()

	return &DatabaseOptions{
		location: path.Join(wd, settings.Location),
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}

func (o *DatabaseOptions) SetLocation(location string) *DatabaseOptions {
	o.location = location

	return o
}
