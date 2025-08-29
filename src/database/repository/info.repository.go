package repository

import (
	"database/sql"
	"fmt"
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

func (r *Info) CreateOne(info *entity.Info) error {
	res, err := r.db.Exec("INSERT INTO info(name, value, value_type) VALUES(?,?,?);", info.Key, info.ValueToString(), info.ValueType)
	if err != nil {
		logger.Warning(err)

		return err
	}

	ins, err := res.RowsAffected()
	if err != nil || ins == 0 {
		return err
	}

	return nil
}

func (r *Info) UpdateOne(info *entity.Info) error {
	res, err := r.db.Exec("UPDATE info SET value=?, value_type=? WHERE name=?", info.ValueToString(), info.ValueType, info.Key)
	if err != nil {
		return err
	}

	if affected, err := res.RowsAffected(); err != nil {
		return err
	} else if affected != 1 {
		return errors.ErrUnableToUpdate
	}

	return nil
}

func (r *Info) FindOneByKey(key string) (*entity.Info, error) {
	row := r.db.QueryRow("SELECT * FROM info WHERE name = ?;", key)
	if row == nil {
		logger.Warning("Selected row from info table was 'nil'")

		return nil, errors.ErrNotFound
	}

	info := &entity.Info{}

	if err := row.Scan(&info.Key, &info.Value, &info.ValueType); err != nil {
		return nil, err
	}

	return info, nil
}

func (r *Info) FindManyByKeys(keys ...string) ([]entity.Info, error) {
	qsm, args := inClause(keys)

	row, err := r.db.Query(fmt.Sprintf("SELECT * FROM info WHERE name IN (%s)", qsm), args...)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	infos := make([]entity.Info, 0)

	for row.Next() {
		info := entity.Info{}

		if err := row.Scan(&info.Key, &info.Value, &info.ValueType); err != nil {
			logger.Warning("FindManyByPrefix entity parse failed:", err)
		} else {
			infos = append(infos, info)
		}
	}

	return infos, nil
}

func (r *Info) FindManyByPrefix(keyPrefix string) ([]entity.Info, error) {
	row, err := r.db.Query("SELECT * FROM info WHERE name LIKE ? || '%'", keyPrefix)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	infos := make([]entity.Info, 0)

	for row.Next() {
		info := entity.Info{}

		if err := row.Scan(&info.Key, &info.Value, &info.ValueType); err != nil {
			logger.Warning("FindManyByPrefix entity parse failed:", err)
		} else {
			infos = append(infos, info)
		}
	}

	return infos, nil
}
