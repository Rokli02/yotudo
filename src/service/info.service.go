package service

import (
	"fmt"
	"yotudo/src/database/entity"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
)

type InfoService struct {
	infoRepository *repository.Info
}

func NewInfoService(infoRepository *repository.Info) *InfoService {
	return &InfoService{
		infoRepository: infoRepository,
	}
}

const (
	m_WINDOW_WIDTH_KEY  = "window_width"
	m_WINDOW_HEIGHT_KEY = "window_height"
)

func (s *InfoService) GetWindowSize() (int, int) {
	width := 1280
	height := 768

	infos, err := s.infoRepository.FindManyByKeys(m_WINDOW_WIDTH_KEY, m_WINDOW_HEIGHT_KEY)
	if err != nil {
		logger.Error(err)
		return width, height
	}

	for _, info := range infos {
		switch info.Key {
		case m_WINDOW_WIDTH_KEY:
			{
				if t_width, err := info.GetValue(); err == nil {
					width = t_width.(int)
				}
			}
		case m_WINDOW_HEIGHT_KEY:
			{
				if t_height, err := info.GetValue(); err == nil {
					height = t_height.(int)
				}
			}
		}
	}

	return width, height
}

func (s *InfoService) SetWindowSize(width, height int) error {
	if width < 300 || height < 200 {
		return fmt.Errorf("window size must be bigger than (300,200) but got (%d,%d)", width, height)
	}

	if err := s.infoRepository.UpdateOne(&entity.Info{Key: m_WINDOW_WIDTH_KEY, Value: width, ValueType: entity.IntValue}); err != nil {
		return err
	}

	if err := s.infoRepository.UpdateOne(&entity.Info{Key: m_WINDOW_HEIGHT_KEY, Value: height, ValueType: entity.IntValue}); err != nil {
		return err
	}

	return nil
}
