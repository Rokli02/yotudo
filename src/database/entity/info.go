package entity

import (
	"strconv"
	"yotudo/src/lib/logger"
)

type Info struct {
	Id        int64
	Key       string
	Value     any
	ValueType InfoType
}

var _ Entity = Info{}

type InfoType uint8

const (
	StringValue InfoType = iota
	IntValue
	DoubleValue
	BoolValue
)

func (Info) Template() string {
	return `
	CREATE TABLE info (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		value TEXT NOT NULL,
		value_type TINYINT DEFAULT 0
	);

	CREATE INDEX info_name_index ON info(name);
	`
}

func (i Info) Migration(currentVersion [3]int) []Migration {
	migrations := []Migration{
		{
			Version:   [3]int{0, 1, 0},
			Migration: i.Template(),
		},
	}

	return migrationsOfVersion(migrations, currentVersion)
}

func (i *Info) GetValue() any {
	switch iVal := i.Value.(type) {
	case string:
		switch i.ValueType {
		case StringValue:
			return iVal
		case BoolValue:
			return iVal != "0"
		case IntValue:
			res, err := strconv.Atoi(iVal)
			if err != nil {
				logger.Warning(err)

				return 0
			}

			return res
		case DoubleValue:
			res, err := strconv.ParseFloat(iVal, 64)
			if err != nil {
				logger.Warning(err)

				return 0.0
			}

			return res
		}
	default:
		return i.Value
	}

	logger.Error("Couldn't parse INFO value", i.Key, i.Value)

	return nil
}
