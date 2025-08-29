package entity

import (
	"fmt"
	"strconv"
	"yotudo/src/lib/logger"
	"yotudo/src/settings"
)

type Info struct {
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
	return fmt.Sprintf(`
	CREATE TABLE info (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		value TEXT NOT NULL,
		value_type TINYINT DEFAULT 0
	);

	CREATE INDEX info_name_index ON info(name);
	INSERT INTO info(name, value, value_type) VALUES('version', '%s', %d);
	`, settings.Global.Database.Version, StringValue)
}

func (i Info) Migration(currentVersion MigrationVersion) []Migration {
	migrations := []Migration{
		{
			Version: MigrationVersion{1, 0, 0},
			Migration: `
				ALTER TABLE info RENAME TO info_old;

				CREATE TABLE info (
					name TEXT PRIMARY KEY,
					value TEXT NOT NULL,
					value_type TINYINT DEFAULT 0
				);

				INSERT INTO info(name, value, value_type) SELECT name, value, value_type FROM info_old;

				DROP TABLE info_old;`,
		},
		{
			Version: MigrationVersion{1, 0, 1},
			Migration: fmt.Sprintf(`
				INSERT INTO info(name, value, value_type) VALUES('window_width', '1280', %d), ('window_height', '768', %d);
			`, IntValue, IntValue),
		},
	}

	return MigrationsByVersion(migrations, currentVersion)
}

func (i *Info) ValueToString() string {
	switch iVal := i.Value.(type) {
	case string:
		return iVal
	default:
		switch i.ValueType {
		case StringValue:
			return iVal.(string)
		case BoolValue:
			if iVal.(bool) {
				return "1"
			}

			return "0"
		case IntValue:
			return fmt.Sprintf("%d", iVal.(int))
		case DoubleValue:
			return fmt.Sprintf("%f", iVal.(float64))
		}
	}

	logger.Error("SHOULD_NOT_REACH end of 'ValueToString' function")

	panic("SHOULD_NOT_REACH")
}

func (i *Info) GetValue() (any, error) {
	switch iVal := i.Value.(type) {
	case string:
		switch i.ValueType {
		case StringValue:
			return iVal, nil
		case BoolValue:
			return iVal != "0", nil
		case IntValue:
			res, err := strconv.Atoi(iVal)
			if err != nil {
				logger.Warning(err)

				return 0, err
			}

			return res, nil
		case DoubleValue:
			res, err := strconv.ParseFloat(iVal, 64)
			if err != nil {
				logger.Warning(err)

				return 0.0, err
			}

			return res, nil
		}
	default:
		return i.Value, nil
	}

	logger.Error("SHOULD_NOT_REACH end of 'GetValue' function")

	panic("SHOULD_NOT_REACH")
}
