package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"yotudo/src/database/entity"
	"yotudo/src/database/errors"
	"yotudo/src/lib/logger"
)

type Contributor struct {
	db *sql.DB
}

func NewContributorRepository(db *sql.DB) *Contributor {
	return &Contributor{db: db}
}

func (c *Contributor) FindByMusicId(conn Connection, musicId int64) []entity.Author {
	if conn == nil {
		conn = c.db
	}

	rows, err := conn.Query(`
		SELECT contributor.author_id, author.name FROM contributor 
		JOIN author ON author.id = contributor.author_id 
		WHERE contributor.music_id = ?`,
		musicId,
	)
	if err != nil {
		logger.Error(err)

		return []entity.Author{}
	}
	defer rows.Close()

	authors := make([]entity.Author, 0)

	for rows.Next() {
		author := entity.Author{}

		if err := rows.Scan(&author.Id, &author.Name); err != nil {
			logger.Warning("Contributor.FindByMusicId:", err)
		} else {
			authors = append(authors, author)
		}
	}

	return authors
}

func (c *Contributor) SaveMany(conn Connection, musicId int64, authorIds []int64) (int64, error) {
	if musicId <= 0 || len(authorIds) == 0 {
		return 0, nil
	}

	if conn == nil {
		conn = c.db
	}

	args := make([]any, len(authorIds))
	values := make([]string, len(authorIds))
	for i, authorId := range authorIds {
		values[i] = fmt.Sprintf("(%d, ?)", musicId)
		args[i] = authorId
	}

	query := fmt.Sprintf("INSERT INTO contributor (music_id, author_id) VALUES%s;", strings.Join(values, ", "))
	res, err := c.db.Exec(query, args...)
	if err != nil {
		logger.Error("Contributor.SaveMany:", err)

		return 0, errors.ErrUnableToSave
	}

	if inserted, err := res.RowsAffected(); err != nil {
		logger.Warning(err)
	} else if inserted != int64(len(authorIds)) {
		logger.WarningF("Tried to insert %d records into \"contributor\" table, but only %d was inserted succesfully", len(authorIds), inserted)
	} else {
		return inserted, nil
	}

	return 0, errors.ErrUnknown
}

func (c *Contributor) DeleteMany(conn Connection, musicId int64, authorIds []int64) (int64, error) {
	if musicId <= 0 || len(authorIds) == 0 {
		return 0, nil
	}

	if conn == nil {
		conn = c.db
	}

	qms, args := inClause(authorIds, musicId)

	query := fmt.Sprintf("DELETE FROM contributor WHERE music_id=? AND author_id IN(%s)", qms)
	res, err := c.db.Exec(query, args...)
	if err != nil {
		logger.Error("Contributor.DeleteMany:", err)

		return 0, errors.ErrUnableToDelete
	}

	if deleted, err := res.RowsAffected(); err != nil {
		logger.Warning(err)
	} else if deleted != int64(len(authorIds)) {
		logger.WarningF("Tried to delete %d records into \"contributor\" table, but only %d was removed succesfully", len(authorIds), deleted)
	} else {
		return deleted, nil
	}

	return 0, errors.ErrUnknown
}

func (c *Contributor) DeleteManyByMusicId(conn Connection, musicId int64) error {
	if conn == nil {
		conn = c.db
	}

	if _, err := conn.Exec("DELETE FROM contributor WHERE music_id=?", musicId); err != nil {
		logger.Error(err)

		return errors.ErrUnableToDelete
	}

	return nil
}
