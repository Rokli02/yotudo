package repository

import (
	"fmt"
	"yotudo/src/database"
	"yotudo/src/database/entity"
	"yotudo/src/database/errors"
	"yotudo/src/lib/logger"
)

type InfoRepository struct{}

var GlobalInfoRepository *InfoRepository = nil

func (i *InfoRepository) CreateOne(info *entity.Info) error {
	res, err := database.Instance.Exec("INSERT INTO info(name, value, value_type) VALUES(?,?,?);", info.Key, info.ValueToString(), info.ValueType)
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

func (i *InfoRepository) UpdateOne(info *entity.Info) error {
	res, err := database.Instance.Exec("UPDATE info SET value=?, value_type=? WHERE name=?", info.ValueToString(), info.ValueType, info.Key)
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

func (i *InfoRepository) FindOneByKey(key string) (*entity.Info, error) {
	row := database.Instance.QueryRow("SELECT name, value, value_type FROM info WHERE name = ?;", key)
	if row == nil {
		logger.Warning("Selected row from info table was 'nil'")

		return nil, errors.ErrNotFound
	}

	info := &entity.Info{}

	if err := info.FromScan(row.Scan); err != nil {
		return nil, err
	}

	return info, nil
}

func (i *InfoRepository) FindManyByKeys(keys ...string) ([]entity.Info, error) {
	qsm, args := inClause(keys)

	row, err := database.Instance.Query(fmt.Sprintf("SELECT name, value, value_type FROM info WHERE name IN (%s)", qsm), args...)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	infos := make([]entity.Info, 0)

	for row.Next() {
		info := entity.Info{}

		if err := info.FromScan(row.Scan); err != nil {
			logger.Warning("FindManyByPrefix entity parse failed:", err)
		} else {
			infos = append(infos, info)
		}
	}

	return infos, nil
}

func (i *InfoRepository) FindManyByPrefix(keyPrefix string) ([]entity.Info, error) {
	row, err := database.Instance.Query("SELECT name, value, value_type FROM info WHERE name LIKE ? || '%'", keyPrefix)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	infos := make([]entity.Info, 0)

	for row.Next() {
		info := entity.Info{}

		if err := info.FromScan(row.Scan); err != nil {
			logger.Warning("FindManyByPrefix entity parse failed:", err)
		} else {
			infos = append(infos, info)
		}
	}

	return infos, nil
}
