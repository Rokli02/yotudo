package entity_test

import (
	"testing"
	"yotudo/src/database/entity"
)

func TestMigrationVersionSetFromText(t *testing.T) {
	mv := entity.MigrationVersion{}
	expectedVersion := "0.1.15"
	mv.SetFromText(expectedVersion)

	if mv[0] != 0 || mv[1] != 1 || mv[2] != 15 {
		t.Errorf("Error during version parse from string to int slice, expected [%s], got [%d.%d.%d]", expectedVersion, mv[0], mv[1], mv[2])
		return
	}
}

func TestMigrationVersionAfterOrEqual(t *testing.T) {
	mv1 := entity.MigrationVersion{0, 5, 2}
	mv2 := entity.MigrationVersion{0, 3, 2}

	if !mv1.EqualOrAfter(mv2) {
		t.Error("[0, 5, 2] should be equal or after [0, 3, 2], but bruh, maybe next time")
		return
	}
}

func TestMigrationVersionAfterOrEqual2(t *testing.T) {
	mv1 := entity.MigrationVersion{0, 5, 2}
	mv2 := entity.MigrationVersion{0, 5, 2}

	if !mv1.EqualOrAfter(mv2) {
		t.Error("[0, 5, 2] should be equal or after [0, 5, 2], but bruh, maybe next time")
		return
	}
}

func TestMigrationVersionAfterOrEqualFail(t *testing.T) {
	mv1 := entity.MigrationVersion{0, 5, 2}
	mv2 := entity.MigrationVersion{0, 6, 2}

	if mv1.EqualOrAfter(mv2) {
		t.Error("[0, 5, 2] should be before [0, 6, 2], but bruh, maybe next time")
		return
	}
}

func TestMigrationGetAll(t *testing.T) {
	migrations := []entity.Migration{
		{
			Version:   entity.MigrationVersion{0, 1, 0},
			Migration: "Numbero uno",
		},
	}
	currentVersion := entity.MigrationVersion{0, 0, 0}
	expectedLength := len(migrations)

	receivedMigrations := entity.MigrationsByVersion(migrations, currentVersion)

	if len(receivedMigrations) != expectedLength {
		t.Errorf("Migration sizes differ, expected to be '%d', but got '%d'", expectedLength, len(receivedMigrations))
		return
	}
}

func TestMigrationGetOnlySome(t *testing.T) {
	migrations := []entity.Migration{
		{
			Version:   entity.MigrationVersion{0, 1, 0},
			Migration: "Numbero uno",
		},
		{
			Version:   entity.MigrationVersion{0, 1, 5},
			Migration: "Dos Tacos",
		},
		{
			Version:   entity.MigrationVersion{0, 1, 8},
			Migration: "Dos Tacos",
		},
	}
	currentVersion := entity.MigrationVersion{0, 1, 0}
	expectedLength := 2

	receivedMigrations := entity.MigrationsByVersion(migrations, currentVersion)

	if len(receivedMigrations) != expectedLength {
		t.Errorf("Migration size expected to be '%d', but got '%d'", expectedLength, len(receivedMigrations))
		return
	}
}

func TestMigrationGetNone(t *testing.T) {
	migrations := []entity.Migration{
		{
			Version:   entity.MigrationVersion{0, 1, 0},
			Migration: "Numbero uno",
		},
		{
			Version:   entity.MigrationVersion{0, 1, 5},
			Migration: "Dos Tacos",
		},
		{
			Version:   entity.MigrationVersion{0, 6, 2},
			Migration: "Dos Tacos",
		},
	}
	currentVersion := entity.MigrationVersion{1, 0, 0}

	receivedMigrations := entity.MigrationsByVersion(migrations, currentVersion)

	if len(receivedMigrations) != 0 {
		t.Errorf("Migration size expected to be '0', but got '%d'", len(receivedMigrations))
		return
	}
}
