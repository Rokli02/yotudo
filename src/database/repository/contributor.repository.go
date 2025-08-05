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

func (c *Contributor) FindByMusicId(musicId int64) []entity.Author {
	rows, err := c.db.Query(`
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
			logger.Warning(err)
		} else {
			authors = append(authors, author)
		}
	}

	return authors
}

func (c *Contributor) SaveMany(musicId int64, authorIds []int64) (int64, error) {
	if musicId <= 0 || len(authorIds) == 0 {
		return 0, errors.ErrNotReceivedInputs
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
		logger.Error(err)

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
