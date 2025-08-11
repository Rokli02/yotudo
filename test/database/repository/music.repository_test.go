package repository_test

import (
	"testing"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
)

func TestSaveMusic(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	authorRepository := repository.NewAuthorRepository(db.Conn)
	musicRepository := repository.NewMusicRepository(db.Conn)

	johnLenon, _ := authorRepository.SaveOne("John Lenon")
	musicId, err := musicRepository.SaveOne(&model.NewMusic{
		Name:      "Test Muzsika",
		Published: 2001,
		Url:       "http://jurta.hu?v=12345",
		Author:    model.OptionalAuthor{Id: &johnLenon.Id},
		GenreId:   1,
	})

	if err != nil {
		t.Error(err)
	}

	row := db.Conn.QueryRow("SELECT * FROM music WHERE id=?", musicId)
	res := make([]any, 11)
	err = row.Scan(&res[0], &res[1], &res[2], &res[3], &res[4], &res[5], &res[6], &res[7], &res[8], &res[9], &res[10])
	if err != nil {
		t.Error(err)

		return
	}

	logger.Debug(res)
}

func TestSaveMusicWithContributors(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	authorRepository := repository.NewAuthorRepository(db.Conn)
	musicRepository := repository.NewMusicRepository(db.Conn)

	johnLenon, _ := authorRepository.SaveOne("John Lenon")
	eltonBro, _ := authorRepository.SaveOne("Elton Bro")
	skibidiGuy, _ := authorRepository.SaveOne("Skibidi Guy")
	musicId, err := musicRepository.SaveOne(&model.NewMusic{
		Name:      "Test Muzsika",
		Published: 2001,
		Url:       "http://jurta.hu?v=12345",
		Author:    model.OptionalAuthor{Id: &johnLenon.Id},
		Contributors: []model.OptionalAuthor{
			{Id: &eltonBro.Id},
			{Id: &skibidiGuy.Id},
		},
		GenreId: 1,
	})

	if err != nil {
		t.Error(err)

		return
	}

	row := db.Conn.QueryRow("SELECT * FROM music WHERE id=?", musicId)
	res := make([]any, 11)
	err = row.Scan(&res[0], &res[1], &res[2], &res[3], &res[4], &res[5], &res[6], &res[7], &res[8], &res[9], &res[10])
	if err != nil {
		t.Error(err)

		return
	}

	contributorIds := make([]int64, 0, 2)
	rows, _ := db.Conn.Query("SELECT author_id FROM contributor WHERE music_id=?;", musicId)
	for rows.Next() {
		var contributorId int64
		rows.Scan(&contributorId)
		contributorIds = append(contributorIds, contributorId)
	}
	rows.Close()

	if len(contributorIds) != 2 {
		t.Errorf("Did not found every added contributors (expected=2, got=%d)", len(contributorIds))
	}

	logger.Debug(res)
}

func TestFindMusicById(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	authorRepository := repository.NewAuthorRepository(db.Conn)
	musicRepository := repository.NewMusicRepository(db.Conn)

	johnLenon, _ := authorRepository.SaveOne("John Lenon")
	musicId, err := musicRepository.SaveOne(&model.NewMusic{
		Name:    "Test Muzsika",
		Album:   "Mi lenne album",
		Url:     "http://jurta.hu?v=12345",
		Author:  model.OptionalAuthor{Id: &johnLenon.Id},
		GenreId: 1,
	})

	if err != nil {
		t.Error(err)
	}

	music, err := musicRepository.FindById(musicId)

	if err != nil {
		t.Error(err)
	} else if music == nil {
		t.Fail()
	}
	logger.Debug(music)
}

func TestFindMusicByIdAfterSavingContributors(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	authorRepository := repository.NewAuthorRepository(db.Conn)
	musicRepository := repository.NewMusicRepository(db.Conn)

	johnLenon, _ := authorRepository.SaveOne("John Lenon")
	eltonJohn, _ := authorRepository.SaveOne("Elton John")
	billyBobber, _ := authorRepository.SaveOne("Billy Bobber")
	musicId, err := musicRepository.SaveOne(&model.NewMusic{
		Name:   "Test Muzsika",
		Album:  "Mi lenne album",
		Url:    "http://jurta.hu?v=12345",
		Author: model.OptionalAuthor{Id: &johnLenon.Id},
		Contributors: []model.OptionalAuthor{
			{Id: &eltonJohn.Id},
			{Id: &billyBobber.Id},
		},
		GenreId: 1,
	})

	if err != nil {
		t.Error(err)
	}

	music, err := musicRepository.FindById(musicId)
	if err != nil {
		t.Error(err)
	}
	if len(music.Contributors) != 2 {
		t.Errorf("Did not found every added contributors (expected=2, got=%d)", len(music.Contributors))
	}

	logger.Debug(music)
}

func TestUpdateOneMusic(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	authorRepository := repository.NewAuthorRepository(db.Conn)
	musicRepository := repository.NewMusicRepository(db.Conn)

	authors, _ := authorRepository.SaveMany([]string{"Test1", "Test2", "Test12", "Test30", "Test23"})
	musicId, err := musicRepository.SaveOne(&model.NewMusic{
		Name:      "Test Muzsika",
		Album:     "Mi lenne album",
		Published: 2000,
		Url:       "http://jurta.hu?v=12345",
		Author:    model.OptionalAuthor{Id: &authors[0].Id},
		Contributors: []model.OptionalAuthor{
			{Id: &authors[2].Id},
			{Id: &authors[4].Id},
		},
		GenreId: 1,
	})

	if err != nil {
		t.Error(err)
	}

	updatedMusic, err := musicRepository.UpdateOne(musicId, &model.UpdateMusic{
		Name:      "Test Módosult Muzsika",
		Published: 2021,
		Url:       "http://jurta.hu?v=12345",
		Author:    model.OptionalAuthor{Id: &authors[0].Id},
		GenreId:   1,
		Status:    1,
		Contributors: []model.OptionalAuthor{
			{Id: &authors[2].Id},
			{Id: &authors[3].Id},
		},
	})

	if err != nil {
		t.Error(err)
	}

	var foundCorrectId uint8 = 0
	for _, contributor := range updatedMusic.Contributors {
		switch contributor.Id {
		case authors[2].Id:
			fallthrough
		case authors[3].Id:
			foundCorrectId++
		case authors[4].Id:
			t.Error("Author should have been deletet from contributors:", contributor)
		}
	}

	if foundCorrectId != 2 {
		t.Errorf("Did not found every added contributors (expected=2, got=%d)", foundCorrectId)
	}

	logger.Debug(updatedMusic)
}

func TestUpdateOneMusic_ErrNotFound(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	authorRepository := repository.NewAuthorRepository(db.Conn)
	musicRepository := repository.NewMusicRepository(db.Conn)

	author, _ := authorRepository.SaveOne("Test1")
	musicId, err := musicRepository.SaveOne(&model.NewMusic{
		Name:      "Test Muzsika",
		Album:     "Mi lenne album",
		Published: 2000,
		Url:       "http://jurta.hu?v=12345",
		Author:    model.OptionalAuthor{Id: &author.Id},
		GenreId:   1,
	})
	if err != nil {
		t.Error(err)
	}

	updatedMusic, err := musicRepository.UpdateOne(musicId+2, &model.UpdateMusic{
		Name:      "Test Módosult Muzsika",
		Published: 2021,
		Url:       "http://jurta.hu?v=12345",
		Author:    model.OptionalAuthor{Id: &author.Id},
		GenreId:   1,
		Status:    1,
	})

	if err == nil {
		t.Error("music entity was found, but it shouldn't have")
	}

	logger.Debug(updatedMusic)
}

func TestFindManyMusic(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	authorRepository := repository.NewAuthorRepository(db.Conn)
	musicRepository := repository.NewMusicRepository(db.Conn)

	logger.Info("Repos created")

	authors, err := authorRepository.SaveMany([]string{"Test1", "Test2", "Test12", "Test30", "Test23"})
	if err != nil {
		t.Error(err)

		return
	}
	_, err = musicRepository.SaveOne(&model.NewMusic{
		Name:      "Test Muzsika",
		Album:     "Mi lenne album",
		Published: 2000,
		Url:       "http://jurta.hu?v=12345",
		Author:    model.OptionalAuthor{Id: &authors[0].Id},
		Contributors: []model.OptionalAuthor{
			{Id: &authors[2].Id},
			{Id: &authors[4].Id},
		},
		GenreId: 1,
	})
	if err != nil {
		t.Error(err)
	}

	_, err = musicRepository.SaveOne(&model.NewMusic{
		Name:      "Komoly Muzsika",
		Published: 1997,
		Url:       "http://jurta.hu?v=86427",
		Author:    model.OptionalAuthor{Id: &authors[1].Id},
		Contributors: []model.OptionalAuthor{
			{Id: &authors[0].Id},
		},
		GenreId: 3,
	})
	if err != nil {
		t.Error(err)
	}

	allMusic, totalCount := musicRepository.FindByPageAndStatus(-1, "Muzs", &model.Page{Size: 0}, []model.Sort{})

	logger.InfoF("Music pagination query ran and found %d records in total", totalCount)

	for _, music := range allMusic {
		logger.Debug(music)
	}
}
