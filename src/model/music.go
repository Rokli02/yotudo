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

type OptionalAuthorGetter interface{ GetOptionalAuthor() *OptionalAuthor }
type OptionalContributorsAccessor interface {
	GetOptionalContributors() []OptionalAuthor
	SetOptionalContributors([]OptionalAuthor)
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
	PicFilename  string
}

func (m *UpdateMusic) GetOptionalAuthor() *OptionalAuthor         { return &m.Author }
func (m *UpdateMusic) GetOptionalContributors() []OptionalAuthor  { return m.Contributors }
func (m *UpdateMusic) SetOptionalContributors(c []OptionalAuthor) { m.Contributors = c }

func (m *UpdateMusic) GetOptionalParams() (Published *int, Album, PicFilename *string) {
	if m.Published != 0 {
		Published = &m.Published
	}

	if m.Album != "" {
		Album = &m.Album
	}

	if m.PicFilename != "" {
		PicFilename = &m.PicFilename
	}

	return
}
