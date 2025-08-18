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
		logger.Error(err)

		return []entity.Genre{}
	}
	defer rows.Close()

	genres = make([]entity.Genre, 0, 1)

	for rows.Next() {
		genre := entity.Genre{}

		err := rows.Scan(&genre.Id, &genre.Name)
		if err != nil {
			logger.Warning("Genre.FindAll:", err)
		} else {
			genres = append(genres, genre)
		}
	}

	return genres
}

func (g *Genre) IsAlreadyUsed(id int64) bool {
	var musicId int64
	row := g.db.QueryRow("SELECT id FROM music WHERE genre_id=? LIMIT 1", id)
	if err := row.Scan(&musicId); err != nil {
		return false
	}

	if musicId != 0 {
		return true
	}

	return false
}

func (g *Genre) SaveOne(name string) (*entity.Genre, error) {
	res, err := g.db.Exec("INSERT INTO genre (name) VALUES (?);", name)
	if err != nil {
		logger.Error("Genre.SaveOne:", err)

		return nil, errors.ErrUnableToSave
	}

	logger.InfoF("Inserted \"%s\" into \"genre\" table", name)

	if id, err := res.LastInsertId(); err != nil {
		logger.Warning(err)

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
		logger.Error("Genre.Rename", err)

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

func (g *Genre) DeleteOne(id int64) error {
	res, err := g.db.Exec("DELETE FROM genre WHERE id=?", id)
	if err != nil {
		logger.Error(err)

		return errors.ErrUnableToDelete
	}

	if affected, err := res.RowsAffected(); err != nil {
		logger.Error(err)

		return errors.ErrUnknown
	} else if affected != 1 {
		return errors.ErrUnableToDelete
	}

	return nil
}
