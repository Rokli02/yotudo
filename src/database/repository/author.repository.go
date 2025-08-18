package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"yotudo/src/database/builders"
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

func (a *Author) FindByPage(filter string, page *model.Page, sort []model.Sort) ([]model.Author, int) {
	args := make([]any, 0)

	totalCountQuery := builders.
		NewQueryBuilder("SELECT COUNT(1) FROM author", &args).
		WithFilter("name", filter).
		WithoutSemicolon().
		Build()

	query := builders.
		NewQueryBuilder(fmt.Sprintf("SELECT id, name, (%s) as total_count FROM author", totalCountQuery), &args).
		WithFilter("name", filter).
		WithSort(sort).
		WithPagination(page).
		Build()

	rows, err := a.db.Query(query, args...)
	if err != nil {
		logger.Error(err)

		return []model.Author{}, 0
	}

	defer rows.Close()

	authors := make([]model.Author, 0, page.Size)

	var totalCount int
	for rows.Next() {
		author := model.Author{}

		err = rows.Scan(&author.Id, &author.Name, &totalCount)
		if err != nil {
			logger.Warning("Author.FindByPage:", err)
		} else {
			authors = append(authors, author)
		}
	}

	return authors, totalCount
}

func (a *Author) SaveOne(name string) (*model.Author, error) {
	var newAuthor *model.Author

	res, err := a.db.Exec("INSERT INTO author (name) VALUES(?);", name)
	if err != nil {
		logger.Error(err)

		return newAuthor, errors.ErrUnableToSave
	}

	if id, err := res.LastInsertId(); err != nil {
		logger.Warning("Author.SaveOne:", err)

		return nil, errors.ErrUnknown
	} else {
		return &model.Author{Id: id, Name: name}, nil
	}
}

func (a *Author) SaveMany(names []string) ([]model.Author, error) {
	if len(names) == 0 {
		return nil, errors.ErrNotReceivedInputs
	}

	tranx, err := a.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := tranx.Rollback(); err != nil && err != sql.ErrTxDone {
			logger.Error("Failed to rollback in \"SaveMany\":", err)
		}
	}()

	stmt, err := tranx.Prepare("INSERT INTO author (name) VALUES (?);")
	if err != nil {
		logger.Error(err)

		return nil, errors.ErrUnknown
	}

	ids := make([]int64, 0, len(names))

	for _, name := range names {
		logger.DebugF("Saving Author(name=%s)", name)

		res, err := stmt.Exec(name)
		if err != nil {
			logger.Error("Author.SaveMany:", err)

			return nil, errors.ErrUnableToSave
		}

		id, err := res.LastInsertId()
		if err != nil {
			logger.Warning(err)

			return nil, errors.ErrUnknown
		}

		ids = append(ids, id)
	}

	if err = stmt.Close(); err != nil {
		logger.Error(err)
	}

	if err = tranx.Commit(); err != nil {
		logger.Error(err)

		return nil, err
	}

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

func (a *Author) IsReferencingToMusic(musicId int64) bool {
	ctx, cancelCtx := context.WithTimeout(context.Background(), time.Second*4)
	defer cancelCtx()

	activeChannelCount := 2
	musicChan := a.isReferencedQuery(ctx, "SELECT COUNT(1) FROM music WHERE author_id=? LIMIT 1", musicId)
	contributorChan := a.isReferencedQuery(ctx, "SELECT COUNT(1) FROM contributor WHERE author_id=? LIMIT 1", musicId)
	isRefed := false

	defer func() {
		close(musicChan)
		close(contributorChan)
	}()

	for activeChannelCount > 0 {
		select {
		case <-ctx.Done():
			logger.Warning("Query timed out, so \"IsReferencingToMusic\" was marked as 'true'")

			return true

		case isReferenced := <-musicChan:
			activeChannelCount--
			isRefed = isRefed || isReferenced

		case isReferenced := <-contributorChan:
			activeChannelCount--
			isRefed = isRefed || isReferenced
		}
	}

	return isRefed
}

func (a *Author) isReferencedQuery(ctx context.Context, queryString string, musicId int64) chan bool {
	isRefChan := make(chan bool)

	go func() {
		var count int

		row := a.db.QueryRowContext(ctx, queryString, musicId)
		if err := row.Scan(&count); err != nil {
			logger.Error(err)

			isRefChan <- true

			return
		}

		isRefChan <- count != 0
	}()

	return isRefChan
}

func (a *Author) DeleteOne(id int64) bool {
	res, err := a.db.Exec("DELETE FROM author WHERE id=?", id)
	if err != nil {
		logger.Error(err)

		return false
	}

	if affected, err := res.RowsAffected(); err != nil {
		logger.Warning("Author.DeleteOne:", err)
	} else if affected == 0 {
		logger.WarningF("Couldn't delete author(%d) for some reason", id)
	} else {
		return true
	}

	return false
}
