package database

import (
	settingsModule "yotudo/src/settings"
)

type DatabaseOptions struct {
	location string
}

func DefaultDatabaseOptions(settings settingsModule.DatabaseSettings) *DatabaseOptions {
	return &DatabaseOptions{
		location: settings.Location,
	}
}

func (o *DatabaseOptions) SetLocation(location string) *DatabaseOptions {
	o.location = location

	return o
}
