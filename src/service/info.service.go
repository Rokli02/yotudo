package service

import "yotudo/src/database/repository"

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

// func (s *InfoService) GetWindowSize() (int, int) {
// 	s.infoRepository.FindManyByKeys(m_WINDOW_WIDTH_KEY, m_WINDOW_HEIGHT_KEY)
// }
