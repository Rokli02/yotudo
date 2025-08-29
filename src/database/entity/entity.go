package entity

import (
	"strconv"
	"strings"
	"yotudo/src/lib/logger"
)

type Entity interface {
	Template() string
	Migration(currentVersion MigrationVersion) []Migration
}

type Migration struct {
	Version   MigrationVersion
	Migration string
}

type MigrationVersion [3]int16

func (v *MigrationVersion) EqualOrAfter(a MigrationVersion) bool {
	var vint int64 = int64(v[0])<<32 | int64(v[1])<<16 | int64(v[2])
	var aint int64 = int64(a[0])<<32 | int64(a[1])<<16 | int64(a[2])

	return vint >= aint
}

func (v *MigrationVersion) SetFromText(version string) {
	versionTexts := strings.Split(version, ".")
	for i := 0; i < len(versionTexts) || i < 3; i++ {
		if vt, err := strconv.Atoi(versionTexts[i]); err == nil {
			v[i] = int16(vt)
		} else {
			logger.Error("Couldn't get version number:", err)
		}
	}
}

func MigrationsByVersion(migrations []Migration, currentVersion MigrationVersion) []Migration {
	for i := len(migrations) - 1; i >= 0; i-- {
		migrationVersion := migrations[i].Version

		if currentVersion.EqualOrAfter(migrationVersion) {
			return migrations[i+1:]
		}
	}

	return migrations
}
