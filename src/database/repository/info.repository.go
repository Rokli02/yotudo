package repository

import (
	"database/sql"
	"yotudo/src/database/entity"
	"yotudo/src/database/errors"
	"yotudo/src/lib/logger"
)

type Info struct {
	db *sql.DB
}

func NewInfoRepository(db *sql.DB) *Info {
	return &Info{db: db}
}

func (r *Info) CreateOne(info *entity.Info) bool {
	res, err := r.db.Exec("INSERT INTO info(name, value, value_type) VALUES(?,?,?);", info.Key, info.Value, info.ValueType)
	if err != nil {
		logger.Warning(err)

		return false
	}

	ins, err := res.RowsAffected()
	if err != nil || ins == 0 {
		return false
	}

	return true
}

func (r *Info) FindOneByKey(key string) (*entity.Info, error) {
	row := r.db.QueryRow("SELECT * FROM info WHERE name = 'version';")
	if row == nil {
		logger.Warning("Selected row from info table was 'nil'")

		return nil, errors.ErrNotFound
	}

	info := &entity.Info{}

	if err := row.Scan(&info.Id, &info.Key, &info.Value, &info.ValueType); err != nil {
		logger.Warning(err)

		info = nil
	}

	return info, nil
}
