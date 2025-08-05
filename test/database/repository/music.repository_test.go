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
	musicId := musicRepository.SaveOne(&model.NewMusic{
		Name:      "Test Muzsika",
		Published: 2001,
		Url:       "http://jurta.hu?v=12345",
		AuthorId:  johnLenon.Id,
		GenreId:   1,
	})

	logger.Debug("MusicId:", musicId)

	row := db.Conn.QueryRow("SELECT * FROM music WHERE id=?", musicId)
	res := make([]any, 11)
	err := row.Scan(&res[0], &res[1], &res[2], &res[3], &res[4], &res[5], &res[6], &res[7], &res[8], &res[9], &res[10])
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	logger.Debug(res)
}

func TestFindMusicById(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	authorRepository := repository.NewAuthorRepository(db.Conn)
	musicRepository := repository.NewMusicRepository(db.Conn)

	johnLenon, _ := authorRepository.SaveOne("John Lenon")
	musicId := musicRepository.SaveOne(&model.NewMusic{
		Name:     "Test Muzsika",
		Album:    "Mi lenne album",
		Url:      "http://jurta.hu?v=12345",
		AuthorId: johnLenon.Id,
		GenreId:  1,
	})

	music := musicRepository.FindById(musicId)
	if music == nil {
		t.Fail()
	}
	logger.Debug(music)
}

func TestFindMusicByIdAfterSavingContributors(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	authorRepository := repository.NewAuthorRepository(db.Conn)
	musicRepository := repository.NewMusicRepository(db.Conn)
	contributorRepository := repository.NewContributorRepository(db.Conn)

	johnLenon, _ := authorRepository.SaveOne("John Lenon")
	eltonJohn, _ := authorRepository.SaveOne("Elton John")
	billyBobber, _ := authorRepository.SaveOne("Billy Bobber")
	musicId := musicRepository.SaveOne(&model.NewMusic{
		Name:     "Test Muzsika",
		Album:    "Mi lenne album",
		Url:      "http://jurta.hu?v=12345",
		AuthorId: johnLenon.Id,
		GenreId:  1,
	})
	savedContributors, err := contributorRepository.SaveMany(musicId, []int64{eltonJohn.Id, billyBobber.Id})

	if err != nil {
		t.Error(err)
	}

	logger.Debug("Saved Contributors:", savedContributors)

	music := musicRepository.FindById(musicId)
	if music == nil {
		t.Fail()
	}
	logger.Debug(music)
}

func TestUpdateOneMusic(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	authorRepository := repository.NewAuthorRepository(db.Conn)
	musicRepository := repository.NewMusicRepository(db.Conn)
	contributorRepository := repository.NewContributorRepository(db.Conn)

	authors, _ := authorRepository.SaveMany([]string{"Test1", "Test2", "Test12", "Test30", "Test23"})
	musicId := musicRepository.SaveOne(&model.NewMusic{
		Name:      "Test Muzsika",
		Album:     "Mi lenne album",
		Published: 2000,
		Url:       "http://jurta.hu?v=12345",
		AuthorId:  authors[0].Id,
		GenreId:   1,
	})
	contributorRepository.SaveMany(musicId, []int64{authors[2].Id, authors[4].Id})
	updatedMusic := musicRepository.UpdateOne(musicId, &model.UpdateMusic{
		Name:           "Test Módosult Muzsika",
		Published:      2021,
		Url:            "http://jurta.hu?v=12345",
		AuthorId:       authors[0].Id,
		GenreId:        1,
		Status:         1,
		ContributorIds: []int64{authors[2].Id, authors[3].Id},
	})

	logger.Debug(updatedMusic)

}

func TestFindManyMusic(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	authorRepository := repository.NewAuthorRepository(db.Conn)
	musicRepository := repository.NewMusicRepository(db.Conn)
	contributorRepository := repository.NewContributorRepository(db.Conn)

	authors, _ := authorRepository.SaveMany([]string{"Test1", "Test2", "Test12", "Test30", "Test23"})
	musicId1 := musicRepository.SaveOne(&model.NewMusic{
		Name:      "Test Muzsika",
		Album:     "Mi lenne album",
		Published: 2000,
		Url:       "http://jurta.hu?v=12345",
		AuthorId:  authors[0].Id,
		GenreId:   1,
	})
	contributorRepository.SaveMany(musicId1, []int64{authors[2].Id, authors[4].Id})

	musicId2 := musicRepository.SaveOne(&model.NewMusic{
		Name:      "Komoly Muzsika",
		Published: 1997,
		Url:       "http://jurta.hu?v=86427",
		AuthorId:  authors[1].Id,
		GenreId:   3,
	})
	contributorRepository.SaveMany(musicId2, []int64{authors[0].Id})

	allMusic := musicRepository.FindByPageAndStatus(-1, "", model.Page{}, []model.Sort{})

	for _, music := range allMusic {
		logger.Debug(music)
	}
}
