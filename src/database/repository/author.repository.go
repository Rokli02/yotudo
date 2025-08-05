package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"yotudo/src/database/errors"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
)

type Author struct {
	db *sql.DB
}

func NewAuthorRepository(db *sql.DB) *Author {
	return &Author{db: db}
}

func (a *Author) FindByPage(filter string, page model.Page, sort []model.Sort) []model.Author {
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString("SELECT id, name FROM author")
	args := make([]any, 0)

	appendQueryWithFilter(filter, &queryBuilder, &args)
	appendQueryWithSort(sort, &queryBuilder)
	appendQueryWithPagination(&page, &queryBuilder, &args)

	queryBuilder.WriteString(";")

	logger.Debug(queryBuilder.String())

	rows, err := a.db.Query(queryBuilder.String(), args...)
	if err != nil {
		logger.Error(err)

		return []model.Author{}
	}

	defer rows.Close()

	authors := make([]model.Author, 0)

	for rows.Next() {
		author := model.Author{}

		err = rows.Scan(&author.Id, &author.Name)
		if err != nil {
			logger.Warning(err)
		} else {
			authors = append(authors, author)
		}
	}

	return authors
}

func (a *Author) SaveOne(name string) (*model.Author, error) {
	var newAuthor *model.Author

	res, err := a.db.Exec("INSERT INTO author (name) VALUES(?);", name)
	if err != nil {
		logger.Warning(err)

		return newAuthor, errors.ErrUnableToSave
	}

	if id, err := res.LastInsertId(); err != nil {
		logger.Warning(err)

		return nil, errors.ErrUnknown
	} else {
		return &model.Author{Id: id, Name: name}, nil
	}
}

func (a *Author) SaveMany(names []string) ([]model.Author, error) {
	if len(names) == 0 {
		return nil, errors.ErrNotReceivedInputs
	}

	stmt, err := a.db.Prepare("INSERT INTO author (name) VALUES (?);")
	if err != nil {
		logger.Error(err)

		return nil, errors.ErrUnknown
	}

	ids := make([]int64, 0, len(names))

	for _, name := range names {
		res, err := stmt.Exec(name)
		if err != nil {
			logger.Warning(err)

			return nil, errors.ErrUnableToSave
		}

		id, err := res.LastInsertId()
		if err != nil {
			logger.Error(err)

			return nil, errors.ErrUnknown
		}

		ids = append(ids, id)
	}

	stmt.Close()

	qms, args := inClause(ids)
	rows, err := a.db.Query(fmt.Sprintf("SELECT id, name FROM author WHERE id IN (%s)", qms), args...)
	if err != nil {
		logger.Error(err)

		return nil, errors.ErrUnableToQuery
	}
	defer rows.Close()

	authors := make([]model.Author, 0, len(ids))
	for rows.Next() {
		author := model.Author{}

		if err := rows.Scan(&author.Id, &author.Name); err != nil {
			logger.Warning(err)

			return nil, errors.ErrUnknown
		} else if author.Id != 0 {
			authors = append(authors, author)
		}
	}

	return authors, nil
}

func (a *Author) DeleteOne(id int64) bool {
	res, err := a.db.Exec("DELETE FROM author WHERE id=?", id)
	if err != nil {
		logger.Warning(err)

		return false
	}

	if affected, err := res.RowsAffected(); err != nil {
		logger.Warning(err)
	} else if affected == 0 {
		logger.WarningF("Couldn't delete author(%d) for some reason", id)
	} else {
		return true
	}

	return false
}
