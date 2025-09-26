package service

import (
	"fmt"
	"yotudo/src/database/entity"
	"yotudo/src/database/repository"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
)

type InfoService struct{}

var GlobalInfoService *InfoService = nil

const (
	m_WINDOW_WIDTH_KEY  = "window_width"
	m_WINDOW_HEIGHT_KEY = "window_height"
)

const (
	m_SERVER_PORT_KEY = "server_port"
	m_SERVER_HOST_KEY = "server_host"
)

func (s *InfoService) GetWindowSize() (int, int) {
	width := 1280
	height := 768

	infos, err := repository.GlobalInfoRepository.FindManyByKeys(m_WINDOW_WIDTH_KEY, m_WINDOW_HEIGHT_KEY)
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

	if err := repository.GlobalInfoRepository.UpdateOne(&entity.Info{Key: m_WINDOW_WIDTH_KEY, Value: width, ValueType: entity.IntValue}); err != nil {
		return err
	}

	if err := repository.GlobalInfoRepository.UpdateOne(&entity.Info{Key: m_WINDOW_HEIGHT_KEY, Value: height, ValueType: entity.IntValue}); err != nil {
		return err
	}

	return nil
}

func (s *InfoService) GetServerConfig() *model.ServerConfig {
	config := model.ServerConfig{}

	infos, err := repository.GlobalInfoRepository.FindManyByKeys(m_SERVER_PORT_KEY, m_SERVER_HOST_KEY)
	if err != nil {
		logger.Error(err)
		return &config
	}

	for _, info := range infos {
		switch info.Key {
		case m_SERVER_PORT_KEY:
			{
				if value, err := info.GetValue(); err == nil {
					config.Port = value.(int)
				}
			}
		case m_SERVER_HOST_KEY:
			{
				if value, err := info.GetValue(); err == nil {
					config.IsHosted = value.(bool)
				}
			}
		}
	}

	return &config
}

func (s *InfoService) SetServerConfig(config *model.ServerConfig) error {
	if err := repository.GlobalInfoRepository.UpdateOne(&entity.Info{
		Key:       m_SERVER_HOST_KEY,
		ValueType: entity.BoolValue,
		Value:     config.IsHosted,
	}); err != nil {
		return err
	}

	if err := repository.GlobalInfoRepository.UpdateOne(&entity.Info{
		Key:       m_SERVER_PORT_KEY,
		ValueType: entity.IntValue,
		Value:     config.Port,
	}); err != nil {
		return err
	}

	return nil
}
