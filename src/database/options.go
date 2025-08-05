package database

import "yotudo/src/settings"

type DatabaseOptions struct {
	location string
}

func DefaultDatabaseOptions() *DatabaseOptions {
	return &DatabaseOptions{
		location: settings.Global.Database.Location,
	}
}

func (o *DatabaseOptions) SetLocation(location string) *DatabaseOptions {
	o.location = location

	return o
}
