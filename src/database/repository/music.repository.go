package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"yotudo/src/database/builders"
	"yotudo/src/database/errors"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
)

type Music struct {
	db                    *sql.DB
	contributorRepository *Contributor
}

func NewMusicRepository(db *sql.DB) *Music {
	return &Music{
		db:                    db,
		contributorRepository: NewContributorRepository(db),
	}
}

func (m *Music) FindByPageAndStatus(status int, filter string, page *model.Page, sort []model.Sort) ([]model.Music, int) {
	args := make([]any, 0)
	totalCountQuery := builders.NewQueryBuilder("SELECT COUNT(1) FROM music", &args).
		WithFilter("name", filter).
		WithCondition("status", status, func(value any) bool {
			return value.(int) > -1
		}).
		WithoutSemicolon().
		Build()

	query := builders.NewQueryBuilder(fmt.Sprintf(
		`SELECT
			m.id, m.name, m.published, m.album, m.url, m.filename, m.pic_filename, m.status,
			a.id, a.name, genre.id, genre.name, ac.id, ac.name, (%s) as total_count
		FROM music AS m 
		JOIN author AS a ON m.author_id = a.id
		JOIN genre ON m.genre_id = genre.id
		LEFT JOIN contributor ON m.id = contributor.music_id
		LEFT JOIN author AS ac ON contributor.author_id = ac.id`,
		totalCountQuery,
	), &args).
		WithFilter("m.name", filter).
		WithCondition("m.status", status, func(value any) bool {
			return value.(int) > -1
		}).
		WithSort(sort).
		WithPagination(page).
		Build()

	rows, err := m.db.Query(query, args...)
	if err != nil {
		logger.Error(err)

		return []model.Music{}, 0
	}

	defer rows.Close()

	musics := make([]model.Music, 0)
	lastIndex := -1

	var totalCount int
	for rows.Next() {
		currentMusic := model.Music{
			Author: model.Author{},
			Genre:  model.Genre{},
		}
		contributor := model.Author{}

		if err := rows.Scan(
			&currentMusic.Id, &currentMusic.Name, &currentMusic.Published, &currentMusic.Album, &currentMusic.Url, &currentMusic.Filename, &currentMusic.PicFilename, &currentMusic.Status,
			&currentMusic.Author.Id, &currentMusic.Author.Name, &currentMusic.Genre.Id, &currentMusic.Genre.Name, &contributor.Id, &contributor.Name, &totalCount,
		); err != nil {
			logger.Warning(err)
		} else {
			if len(musics) == 0 || musics[lastIndex].Id != currentMusic.Id {
				lastIndex++
				musics = append(musics, currentMusic)
			}

			if len(musics[lastIndex].Contributors) == 0 {
				musics[lastIndex].Contributors = make([]model.Author, 0)
			}

			if contributor.Id != 0 {
				musics[lastIndex].Contributors = append(musics[lastIndex].Contributors, contributor)
			}
		}
	}

	return musics, totalCount
}

func (m *Music) FindById(id int64) (*model.Music, error) {
	var music *model.Music

	rows, err := m.db.Query(
		`SELECT
			m.name, m.published, m.album, m.url, m.filename, m.pic_filename, m.status,
			a.id, a.name, genre.id, genre.name, ac.id, ac.name
		FROM music AS m 
		JOIN author AS a ON m.author_id = a.id
		JOIN genre ON m.genre_id = genre.id
		LEFT JOIN contributor ON m.id = contributor.music_id
		LEFT JOIN author AS ac ON contributor.author_id = ac.id
		WHERE m.id=?;`,
		id,
	)
	if err != nil {
		logger.Warning(err)

		return nil, errors.ErrUnknown
	}

	defer rows.Close()

	music = &model.Music{
		Id:           id,
		Author:       model.Author{},
		Contributors: make([]model.Author, 0),
		Genre:        model.Genre{},
	}

	for rows.Next() {
		var contributorId *int64
		var contributorName *string

		if err := rows.Scan(
			&music.Name, &music.Published, &music.Album, &music.Url, &music.Filename, &music.PicFilename, &music.Status,
			&music.Author.Id, &music.Author.Name, &music.Genre.Id, &music.Genre.Name, &contributorId, &contributorName,
		); err != nil {
			logger.Warning(err)

			return nil, errors.ErrNotFound
		} else if contributorId != nil {
			logger.DebugF("Contributor(id=%d, name=%s)", *contributorId, *contributorName)
			music.Contributors = append(music.Contributors, model.Author{Id: *contributorId, Name: *contributorName})
		}
	}

	return music, nil
}

func (m *Music) SaveOne(newMusic *model.NewMusic) (int64, error) {
	var published *int = nil
	var album *string = nil
	if newMusic.Published != 0 {
		published = &newMusic.Published
	}
	if strings.TrimSpace(newMusic.Album) != "" {
		album = &newMusic.Album
	}

	res, err := m.db.Exec(
		"INSERT INTO music (name, published, album, url, author_id, genre_id, updated_at) VALUES(?,?,?,?,?,?,?);",
		newMusic.Name, published, album, newMusic.Url,
		newMusic.Author.Id, newMusic.GenreId, time.Now().Format(DefaultDateFormat),
	)
	if err != nil {
		logger.Warning(err)

		return 0, errors.ErrUnableToSave
	}
	id, err := res.LastInsertId()
	if err != nil {
		logger.Error(err)

		return 0, errors.ErrUnknown
	}

	contributorIds := make([]int64, len(newMusic.Contributors))
	for i := 0; i < len(newMusic.Contributors); i++ {
		contributorIds[i] = *newMusic.Contributors[i].Id
	}

	_, err = m.contributorRepository.SaveMany(id, contributorIds)
	if err != nil {
		logger.Error(err)
	}

	return id, nil
}

/*
Updates Music and its contributor references too.
*/
func (m *Music) UpdateOne(musicId int64, music *model.UpdateMusic) (updateOneResponse *model.Music, returnError error) {
	updateOneResponse = nil
	returnError = nil

	if music == nil {
		returnError = errors.ErrNotReceivedInputs

		return
	}

	row := m.db.QueryRow(
		"SELECT id FROM music WHERE id=? LIMIT 1",
		musicId,
	)

	var idFromDB int64

	if err := row.Scan(&idFromDB); err != nil {
		logger.Warning(err)
		returnError = errors.ErrNotFound

		return
	} else if idFromDB == 0 {
		returnError = errors.ErrNotFound

		return
	}

	// If a value in 'contributorState' is 'false' must be deleted, if 'true' must be added, otherwise do nothing
	contributorState := map[int64]bool{}

	// Get contributors' ids from 'contributor' table
	rows, err := m.db.Query("SELECT author_id FROM contributor WHERE music_id=?", musicId)
	if err != nil {
		logger.Error(err)
	} else {
		defer rows.Close()

		for rows.Next() {
			var contributorId int64

			if err := rows.Scan(&contributorId); err != nil {
				logger.Warning(err)
			} else if contributorId != 0 {
				contributorState[contributorId] = false
			}
		}
	}

	if len(music.Contributors) > 0 {
		for _, contributor := range music.Contributors {
			if _, found := contributorState[*contributor.Id]; found {
				delete(contributorState, *contributor.Id)
			} else {
				contributorState[*contributor.Id] = true
			}
		}
	}

	trans, err := m.db.Begin()
	if err != nil {
		logger.ErrorF("CRITICAL ERROR: Transaction couldn't start for updating music(id=%d)", musicId)
		logger.Error(err)
		returnError = errors.ErrUnknown

		return
	}

	// If 'updateOneResponse' remained nil something went wrong during execution, so do a db rollback
	// Otherwise completed its purpose, so commit changes
	defer func() {
		if updateOneResponse == nil {
			logger.Debug("UpdateOne function ran into some problem during execution")

			if err := trans.Rollback(); err != nil {
				logger.Error(err)
			}
		}
	}()

	// Update 'music' record based on the given properties
	published, album, filename, picFilename := music.GetOptionalParams()
	res, err := trans.Exec(`
		UPDATE music
		SET author_id=?, name=?, published=?, album=?, genre_id=?, url=?, filename=?, pic_filename=?,
			status=?, updated_at=?
		WHERE id=?`,
		*music.Author.Id, music.Name, published, album, music.GenreId, music.Url, filename, picFilename,
		music.Status, time.Now().Format(DefaultDateFormat), musicId,
	)
	if err != nil {
		logger.Error(err)
		returnError = errors.ErrUnableToUpdate

		return
	}

	if affected, err := res.RowsAffected(); affected == 0 {
		logger.WarningF("Couldn't update music(id=%d) because an error occured: %s", musicId, err)
		returnError = errors.ErrUnableToUpdate

		return
	}

	// Update 'contributors'
	// Contributor list has changed
	// Decide which contributors got removed from music and which got added
	if len(contributorState) != 0 {
		contributorsToDelete := make([]int64, 0)
		contributorsToAdd := make([]int64, 0)

		for key, value := range contributorState {
			if !value {
				contributorsToDelete = append(contributorsToDelete, key)
			} else {
				contributorsToAdd = append(contributorsToAdd, key)
			}
		}

		// Delete those contributors which are not relevant anymore
		if len(contributorsToDelete) != 0 {
			logger.Debug("Contributors To Delete:", contributorsToDelete)

			qms, args := inClause(contributorsToDelete, musicId)
			res, err := trans.Exec(fmt.Sprintf("DELETE FROM contributor WHERE music_id=? AND author_id IN (%s)", qms), args...)
			if err != nil {
				logger.Warning(err)
			}

			if affected, err := res.RowsAffected(); affected != int64(len(contributorsToDelete)) {
				logger.Error("Couldn't delete all of the contributors:", err)
				returnError = errors.ErrUnableToUpdate

				return
			}
		}

		// Add those contributors which were not present already
		if len(contributorsToAdd) != 0 {
			logger.Debug("Contributors To Add:", contributorsToAdd)

			args := make([]any, len(contributorsToAdd))
			values := make([]string, len(contributorsToAdd))
			for i, authorId := range contributorsToAdd {
				values[i] = fmt.Sprintf("(%d, ?)", musicId)
				args[i] = authorId
			}

			query := fmt.Sprintf("INSERT INTO contributor (music_id, author_id) VALUES%s;", strings.Join(values, ", "))
			res, err := trans.Exec(query, args...)
			if err != nil {
				logger.Error(err)
				returnError = errors.ErrUnableToUpdate

				return
			}

			if affected, err := res.RowsAffected(); affected != int64(len(contributorsToAdd)) {
				logger.Error("Couldn't insert all of the new contributors:", err)
				returnError = errors.ErrUnableToUpdate

				return
			}
		}
	}

	if err := trans.Commit(); err != nil {
		logger.Error(err)
		returnError = errors.ErrUnknown

		return
	}

	updateOneResponse, returnError = m.FindById(musicId)

	return
}

func (m *Music) UpdateStatus(id int64, status int) error {
	res, err := m.db.Exec("UPDATE music SET status=? WHERE id=?", status, id)
	if err != nil {
		logger.Error(err)

		return errors.ErrUnableToUpdate
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return errors.ErrUnknown
	}
	if affected == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (m *Music) DeleteOne(id int64) (bool, error) {
	tranx, err := m.db.Begin()
	if err != nil {
		logger.Error(err)

		return false, errors.ErrUnknown
	}
	defer func() {
		if err := tranx.Rollback(); err != nil && err != sql.ErrTxDone {
			logger.Error(err)
		}
	}()

	if _, err := tranx.Exec("DELETE FROM contributor WHERE music_id=?", id); err != nil {
		logger.Error(err)

		return false, errors.ErrUnknown
	}

	if res, err := tranx.Exec("DELETE FROM music WHERE id=?", id); err != nil {
		logger.Error(err)

		return false, errors.ErrUnknown
	} else if affected, err := res.RowsAffected(); err != nil {
		logger.Error(err)

		return false, errors.ErrUnknown
	} else if affected == 0 {
		logger.ErrorF("Couldn't delete music with id (%d)", id)

		return false, errors.ErrUnableToDelete
	}

	if err = tranx.Commit(); err != nil {
		logger.Error(err)

		return false, errors.ErrUnknown
	}

	return true, nil
}
