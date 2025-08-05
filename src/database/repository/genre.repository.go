package repository

import (
	"database/sql"
	"yotudo/src/database/entity"
	"yotudo/src/database/errors"
	"yotudo/src/lib/logger"
)

type Genre struct {
	db *sql.DB
}

func NewGenreRepository(db *sql.DB) *Genre {
	return &Genre{db: db}
}

func (g *Genre) FindAll() []entity.Genre {
	var genres []entity.Genre

	rows, err := g.db.Query("SELECT id, name FROM genre;")
	if err != nil {
		logger.Warning(err)

		return []entity.Genre{}
	}
	defer rows.Close()

	genres = make([]entity.Genre, 0, 1)

	for rows.Next() {
		genre := entity.Genre{}

		err := rows.Scan(&genre.Id, &genre.Name)
		if err != nil {
			logger.Warning(err)
		} else {
			genres = append(genres, genre)
		}
	}

	return genres
}

func (g *Genre) SaveOne(name string) (*entity.Genre, error) {
	res, err := g.db.Exec("INSERT INTO genre (name) VALUES (?);", name)
	if err != nil {
		logger.Warning(err)

		return nil, errors.ErrUnableToSave
	}

	logger.InfoF("Inserted \"%s\" into \"genre\" table", name)

	if id, err := res.LastInsertId(); err != nil {
		logger.Warning(err)

		// #region Extra select insiode SaveOne
		// TODO: Delete comments if function still works
		// row := g.db.QueryRow("SELECT id FROM genre WHERE name=? LIMIT 1", name)
		// if row == nil {
		// 	logger.ErrorF("Genre \"%s\" got inserted into database, but was not found later in a query!", name)
		//
		// 	return nil, errors.ErrUnableToSave
		// }
		//
		// genre := &entity.Genre{Name: name}
		//
		// err = row.Scan(&genre.Id)
		// if err != nil {
		// 	logger.Warning(err)
		// }
		// #endregion

		return nil, errors.ErrUnknown
	} else {
		return &entity.Genre{
			Id:   id,
			Name: name,
		}, nil
	}
}

func (g *Genre) Rename(id int64, newName string) (*entity.Genre, error) {
	res, err := g.db.Exec("UPDATE genre SET name=? WHERE id=?", newName, id)
	if err != nil {
		logger.Warning(err)

		return nil, errors.ErrUnableToUpdate
	}
	updatedRows, err := res.RowsAffected()
	if err != nil {
		logger.Warning(err)

		return nil, errors.ErrUnableToUpdate
	}
	if updatedRows == 0 {
		logger.WarningF("Couldn't rename genre to \"%s\"", newName)

		return nil, errors.ErrUnableToUpdate
	}

	return &entity.Genre{
		Id:   id,
		Name: newName,
	}, nil
}
