package repository_test

import (
	"testing"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
)

func TestSaveAuthors(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	repo := repository.NewAuthorRepository(db.Conn)

	author, err := repo.SaveOne("Énekes 1")
	if err != nil {
		t.Error(err)
	}

	logger.Info("Saved Author", author)
}

func TestSaveAuthorsThenGetSome(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	repo := repository.NewAuthorRepository(db.Conn)

	repo.SaveOne("Énekes 1")
	repo.SaveOne("Színész 7")
	repo.SaveOne("Énekes 2")
	repo.SaveOne("Énekes 5")
	repo.SaveOne("Színész 6")
	repo.SaveOne("Énekes 3")
	repo.SaveOne("Énekes 4")

	filter := "Énes"
	page := model.Page{Page: 0, Size: 3}
	sort := []model.Sort{
		{
			Key: "id",
			Dir: -1,
		},
	}
	logger.InfoF("Parameters: filter(\"%s\"), page(%v), sort(%v)", filter, page, sort)

	authors, totalCount := repo.FindByPage(filter, &page, sort)

	logger.Info("All Authors based on parameters", authors, totalCount)
}

func TestSaveManyAuthors(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	repo := repository.NewAuthorRepository(db.Conn)

	savedAuthors, err := repo.SaveMany([]string{"Tester Énekes", "Profi Énekes", "Amatör Zajkeltő"})
	if err != nil {
		t.Error(err)
	}

	if len(savedAuthors) != 3 {
		t.Fail()
	}

	logger.Debug(savedAuthors)
}

func TestSaveManyAuthorsFailUniqueNameRestraint(t *testing.T) {
	db := getInMemoryDB()
	defer db.Close()
	repo := repository.NewAuthorRepository(db.Conn)

	authorNames := []string{"Egyedi Énekes", "Egyedi Énekes"}

	_, err := repo.SaveMany(authorNames)

	if err == nil {
		t.Errorf("Saved %d authors, with the same name", len(authorNames))
	}
}
