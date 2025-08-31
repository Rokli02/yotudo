package model

import "encoding/json"

type Music struct {
	Id           int64
	Name         string
	Published    *int
	Album        *string
	Url          string
	Filename     *string
	PicFilename  *string
	Status       int8
	Genre        Genre
	Author       Author
	Contributors []Author
}

type OptionalAuthorGetter interface{ GetOptionalAuthor() *OptionalAuthor }
type OptionalContributorsAccessor interface {
	GetOptionalContributors() []OptionalAuthor
	SetOptionalContributors([]OptionalAuthor)
}

func (m *Music) ToUpdateMusic() *UpdateMusic {
	updateMusic := &UpdateMusic{
		Id:      m.Id,
		Name:    m.Name,
		Url:     m.Url,
		Status:  m.Status,
		GenreId: m.Genre.Id,
		Author:  OptionalAuthor{Id: &m.Author.Id, Name: &m.Author.Name},
	}

	if m.Published != nil {
		updateMusic.Published = *m.Published
	}

	if m.Album != nil {
		updateMusic.Album = *m.Album
	}

	if m.Filename != nil {
		updateMusic.Filename = *m.Filename
	}

	if m.PicFilename != nil {
		updateMusic.PicFilename = *m.PicFilename
	}

	if len(m.Contributors) != 0 {
		updateMusic.Contributors = make([]OptionalAuthor, len(m.Contributors))

		for i := 0; i < len(updateMusic.Contributors); i++ {
			updateMusic.Contributors[i] = OptionalAuthor{
				Id:   &m.Contributors[i].Id,
				Name: &m.Contributors[i].Name,
			}
		}

	}

	return updateMusic
}

type NewMusic struct {
	Name         string
	Published    int
	Album        string
	Url          string
	Author       OptionalAuthor
	Contributors []OptionalAuthor
	GenreId      int64
	PicFilename  string
	PicType      string
}

func (m *NewMusic) GetOptionalAuthor() *OptionalAuthor         { return &m.Author }
func (m *NewMusic) GetOptionalContributors() []OptionalAuthor  { return m.Contributors }
func (m *NewMusic) SetOptionalContributors(c []OptionalAuthor) { m.Contributors = c }

type UpdateMusic struct {
	Id           int64
	Name         string
	Published    int
	Album        string
	Url          string
	Author       OptionalAuthor
	Contributors []OptionalAuthor
	Status       int8
	GenreId      int64
	Filename     string
	PicFilename  string
	PicType      string
}

func (m *UpdateMusic) GetOptionalAuthor() *OptionalAuthor         { return &m.Author }
func (m *UpdateMusic) GetOptionalContributors() []OptionalAuthor  { return m.Contributors }
func (m *UpdateMusic) SetOptionalContributors(c []OptionalAuthor) { m.Contributors = c }

func (m *UpdateMusic) GetOptionalParams() (Published *int, Album, Filename, PicFilename *string) {
	if m.Published != 0 {
		Published = &m.Published
	}

	if m.Album != "" {
		Album = &m.Album
	}

	if m.PicFilename != "" {
		PicFilename = &m.PicFilename
	}

	if m.Filename != "" {
		Filename = &m.Filename
	}

	return
}
func (m *UpdateMusic) String() string {
	jb, _ := json.Marshal(m)
	return string(jb)
}
