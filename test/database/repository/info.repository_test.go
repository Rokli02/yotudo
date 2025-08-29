package repository_test

import (
	"strconv"
	"testing"
	"yotudo/src/database/entity"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
)

func TestSaveInfo(t *testing.T) {
	db := getInMemoryDB()
	infoRepository := repository.NewInfoRepository(db.Conn)

	if err := infoRepository.CreateOne(&entity.Info{Key: "test_width", Value: "1200", ValueType: entity.IntValue}); err != nil {
		t.Error(err)
		return
	}
}

func TestGetByIdInfo(t *testing.T) {
	db := getInMemoryDB()
	infoRepository := repository.NewInfoRepository(db.Conn)

	const (
		widthKey  = "test_width"
		heightKey = "test_height"
	)
	expectedValue := 1200

	infoRepository.CreateOne(&entity.Info{Key: widthKey, Value: strconv.Itoa(expectedValue), ValueType: entity.IntValue})
	infoRepository.CreateOne(&entity.Info{Key: heightKey, Value: "800", ValueType: entity.IntValue})

	if info, err := infoRepository.FindOneByKey(widthKey); err != nil {
		t.Error(err)
		return
	} else if returnedValue := info.GetValue().(int); returnedValue != expectedValue {
		t.Errorf("Expected value to be '%d', but got '%d'", expectedValue, returnedValue)
		return
	}
}

func TestGetByIdInfo_NotFound(t *testing.T) {
	db := getInMemoryDB()
	infoRepository := repository.NewInfoRepository(db.Conn)

	const (
		widthKey  = "test_width"
		heightKey = "test_height"
	)

	infoRepository.CreateOne(&entity.Info{Key: widthKey, Value: "1200", ValueType: entity.IntValue})
	infoRepository.CreateOne(&entity.Info{Key: heightKey, Value: "800", ValueType: entity.IntValue})

	if _, err := infoRepository.FindOneByKey("xxx_lehetetlen_kulcs_420__"); err == nil {
		t.Error("Shouldn't have found any result, but it did")
		return
	}
}

func TestGetByKeys(t *testing.T) {
	db := getInMemoryDB()
	infoRepository := repository.NewInfoRepository(db.Conn)

	const (
		widthKey  = "test_width"
		heightKey = "test_height"
	)

	infoRepository.CreateOne(&entity.Info{Key: widthKey, Value: "1200", ValueType: entity.IntValue})
	infoRepository.CreateOne(&entity.Info{Key: heightKey, Value: "800", ValueType: entity.IntValue})

	if infos, err := infoRepository.FindManyByKeys(widthKey, heightKey); err != nil {
		t.Error(err)
		return
	} else if len(infos) != 2 {
		t.Errorf("Should have return with '2' results, but got '%d'", len(infos))
		return
	}
}

func TestGetByKeyPrefix(t *testing.T) {
	db := getInMemoryDB()
	infoRepository := repository.NewInfoRepository(db.Conn)

	const (
		widthKey  = "test_width"
		heightKey = "test_height"
	)

	infoRepository.CreateOne(&entity.Info{Key: widthKey, Value: "1200", ValueType: entity.IntValue})
	infoRepository.CreateOne(&entity.Info{Key: heightKey, Value: "800", ValueType: entity.IntValue})

	if infos, err := infoRepository.FindManyByPrefix("test_"); err != nil {
		t.Error(err)
		return
	} else if len(infos) != 2 {
		t.Errorf("Should have return with '2' results, but got '%d'", len(infos))
		return
	}
}

func TestUpdateOne(t *testing.T) {
	db := getInMemoryDB()
	infoRepository := repository.NewInfoRepository(db.Conn)

	const (
		widthKey  = "test_width"
		heightKey = "test_height"
	)

	expectedValue := 728

	infoRepository.CreateOne(&entity.Info{Key: widthKey, Value: "1200", ValueType: entity.IntValue})
	infoRepository.CreateOne(&entity.Info{Key: heightKey, Value: "800", ValueType: entity.IntValue})

	infoRepository.UpdateOne(&entity.Info{Key: heightKey, Value: strconv.Itoa(expectedValue), ValueType: entity.IntValue})

	if infos, err := infoRepository.FindManyByPrefix("test_"); err != nil {
		t.Error(err)
		return
	} else {
		logger.DebugF("Returned with '%d' elements", len(infos))

		for _, info := range infos {
			if returnedValue := info.GetValue(); info.Key == heightKey && returnedValue != expectedValue {
				t.Errorf("Record '%s' should have been updated to '%d' but was %d", info.Key, expectedValue, returnedValue)
			}
		}
	}
}
