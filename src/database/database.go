package database

import (
	"database/sql"
	"strings"
	"yotudo/src/database/entity"
	"yotudo/src/database/repository"
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
	dbOption := DefaultDatabaseOptions(settings.Global.Database)
	for _, optsFunc := range optsFuncs {
		optsFunc(dbOption)
	}

	db := newDatabase(dbOption)

	isExists := db.isAlreadyInitialized()

	if !isExists {
		logger.Info("Initializing Database ...")

		db.init()
		db.migrateDatabase("0.0.0")
	} else {
		infoRepository := repository.NewInfoRepository(db.Conn)

		info, _ := infoRepository.FindOneByKey("version")
		if info != nil && info.Value != settings.Global.Database.Version {
			logger.WarningF("Current database version is %s, but the newest is %s", info.Value, settings.Global.Database.Version)
			logger.Info("Migrating database to the newest version...")

			db.migrateDatabase(info.ValueToString())
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
	logger.Info("Closing Database ...")
	db.Conn.Close()
}

func (db *Database) isAlreadyInitialized() bool {
	var rootPage int

	if err := db.Conn.QueryRow("SELECT rootpage from sqlite_master WHERE type = 'type' AND name = 'info';").Scan(&rootPage); err != nil {
		logger.Warning("In 'isAlreadyInitialized' warning occured:", err)

		return false
	} else if rootPage == 0 {
		return false
	}

	return true
}

func (db *Database) init() {
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
}

func (db *Database) migrateDatabase(versionText string) {
	// Convert version string into fixed array
	version := entity.MigrationVersion{0, 0, 0}
	version.SetFromText(versionText)

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

	migrationString := sb.String()

	if _, err := db.Conn.Exec(migrationString); err != nil {
		logger.Error("Error during database migration (Couldn't execute built command):", err)
	}
}
