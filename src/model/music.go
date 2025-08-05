package model

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

type NewMusic struct {
	Name           string
	Published      int
	Album          string
	Url            string
	AuthorId       int64
	ContributorIds []int64
	GenreId        int64
}

type UpdateMusic struct {
	Name           string
	Published      int
	Album          string
	Url            string
	Filename       string
	PicFilename    string
	Status         int8
	GenreId        int64
	AuthorId       int64
	ContributorIds []int64
}

func (m *UpdateMusic) GetOptionalParams() (Published *int, Album, Filename, PicFilename *string) {
	if m.Published != 0 {
		Published = &m.Published
	}

	if m.Album != "" {
		Album = &m.Album
	}

	if m.Filename != "" {
		Filename = &m.Filename
	}

	if m.PicFilename != "" {
		PicFilename = &m.PicFilename
	}

	return
}
