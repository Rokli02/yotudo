package repository_test

import (
	"testing"
	"yotudo/src/database/entity"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
)

func TestSaveMusic(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	authorRepository := repository.GlobalAuthorRepository
	musicRepository := repository.GlobalMusicRepository

	johnLenon := Must(authorRepository.SaveOne("John Lenon"))
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
	musicEntity := &entity.Music{}

	Must(musicEntity.FromScan(row.Scan))

	logger.Debug(musicEntity)
}

func TestSaveMusicWithContributors(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	authorRepository := repository.GlobalAuthorRepository
	musicRepository := repository.GlobalMusicRepository

	johnLenon := Must(authorRepository.SaveOne("John Lenon"))
	eltonBro := Must(authorRepository.SaveOne("Elton Bro"))
	skibidiGuy := Must(authorRepository.SaveOne("Skibidi Guy"))
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
	musicEntity := &entity.Music{}

	Must(musicEntity.FromScan(row.Scan))

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

	logger.Debug(musicEntity)
}

func TestSaveMusicWithImage(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	authorRepository := repository.GlobalAuthorRepository
	imageRepository := repository.GlobalImageRepository
	musicRepository := repository.GlobalMusicRepository

	johnLenon := Must(authorRepository.SaveOne("John Lenon"))

	musicId, err := musicRepository.SaveOne(&model.NewMusic{
		Name:      "Test Muzsika",
		Published: 2001,
		Url:       "http://jurta.hu?v=12345",
		Author:    model.OptionalAuthor{Id: &johnLenon.Id},
		Image:     &model.PossiblyNewImage{Name: "newImage.png"},
		GenreId:   1,
	})

	if err != nil {
		t.Error(err)

		return
	}

	row := db.Conn.QueryRow("SELECT * FROM music WHERE id=?", musicId)
	musicEntity := &entity.Music{}

	Must(musicEntity.FromScan(row.Scan))

	if image, err := imageRepository.FindById(nil, *musicEntity.ImageId); err != nil {
		t.Error(err)
		return
	} else {
		logger.Debug("Image got saved:", image)
	}
}

func TestFindMusicById(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	authorRepository := repository.GlobalAuthorRepository
	musicRepository := repository.GlobalMusicRepository

	johnLenon := Must(authorRepository.SaveOne("John Lenon"))
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
	authorRepository := repository.GlobalAuthorRepository
	musicRepository := repository.GlobalMusicRepository

	johnLenon := Must(authorRepository.SaveOne("John Lenon"))
	eltonJohn := Must(authorRepository.SaveOne("Elton John"))
	billyBobber := Must(authorRepository.SaveOne("Billy Bobber"))
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
	authorRepository := repository.GlobalAuthorRepository
	musicRepository := repository.GlobalMusicRepository

	authors := Must(authorRepository.SaveMany([]string{"Test1", "Test2", "Test12", "Test30", "Test23"}))
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
	authorRepository := repository.GlobalAuthorRepository
	musicRepository := repository.GlobalMusicRepository

	author := Must(authorRepository.SaveOne("Test1"))
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
	authorRepository := repository.GlobalAuthorRepository
	musicRepository := repository.GlobalMusicRepository

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
