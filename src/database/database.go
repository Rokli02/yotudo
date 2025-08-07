package database

import (
	"database/sql"
	"os"
	"strconv"
	"strings"
	"yotudo/src/database/entity"
	"yotudo/src/lib/logger"
	"yotudo/src/settings"

	_ "modernc.org/sqlite"
)

type Database struct {
	Conn *sql.DB
}

type DatabaseOptionsFunc func(opts *DatabaseOptions)

func (db *Database) databaseTables() []entity.Entity {
	return []entity.Entity{
		entity.Info{},
		entity.Genre{},
		entity.Status{},
		entity.Author{},
		entity.Music{},
		entity.Contributor{},
	}
}

func LoadDatabase(optsFuncs ...DatabaseOptionsFunc) *Database {
	isExists := isDatabaseExists(settings.Global.Database.Location)

	dbOption := DefaultDatabaseOptions()
	for _, optsFunc := range optsFuncs {
		optsFunc(dbOption)
	}

	db := newDatabase(dbOption)

	InitRepositories(db.Conn)

	if !isExists {
		db.init()
	} else {
		info, _ := InfoRepository.FindOneByKey("version")
		if info != nil && info.Value != settings.Global.Database.Version {
			logger.ErrorF("Current database version is %s, but the newest is %s", info.Value, settings.Global.Database.Version)

			db.migrateDatabase(info.Value.(string))
		}
	}

	return db
}

func newDatabase(dbOption *DatabaseOptions) *Database {
	conn, err := sql.Open("sqlite", dbOption.location)
	if err != nil {
		logger.Error(err)
		panic(-1)
	}

	err = conn.Ping()
	if err != nil {
		logger.Error(err)
		panic(-1)
	}

	return &Database{
		Conn: conn,
	}
}

func (db *Database) Close() {
	logger.Info("Closing Database")
	db.Conn.Close()
}

func (db *Database) init() {
	logger.Info("Initializing Database")

	databaseTableList := db.databaseTables()

	sb := strings.Builder{}

	for _, table := range databaseTableList {
		sb.WriteString(table.Template())
	}

	_, err := db.Conn.Exec(sb.String())
	if err != nil {
		logger.Error(err)
		panic(-1)
	}

	logger.Info("Preloading Database with mandatory datas")

	InfoRepository.CreateOne(&entity.Info{Key: "version", Value: settings.Global.Database.Version, ValueType: entity.StringValue})
}

func isDatabaseExists(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	return true
}

func (db *Database) migrateDatabase(versionText string) {
	// Convert version string into fixed array
	versionTexts := strings.Split(versionText, ".")
	version := [3]int{0, 0, 0}

	for i := 0; i < len(versionTexts) || i < 3; i++ {
		if v, err := strconv.Atoi(versionTexts[i]); err == nil {
			version[i] = v
		} else {
			logger.Error("Error during database migration (Couldn't get version number):", err)
		}
	}

	databaseTableList := db.databaseTables()
	sb := strings.Builder{}

	// Merge migrations into one command
	for _, table := range databaseTableList {
		migrations := table.Migration(version)

		if len(migrations) > 0 {
			for _, migration := range migrations {
				sb.WriteString(migration.Migration)
				sb.WriteString("\n")
			}
		}
	}

	if _, err := db.Conn.Exec(sb.String()); err != nil {
		logger.Error("Error during database migration (Couldn't execute built command):", err)
	}
}
