package service

import (
	"fmt"
	"os"
	pathModule "path"
	"regexp"
	"strings"
	"time"
	"yotudo/src/model"
)

type FileService uint8

var (
	validFilenameRegexp       = regexp.MustCompile(`^([^\\\:\*\?\"\/\<\>\|\0]{1,255})$`)
	invalidFilenameCharacters = regexp.MustCompile(`[\\\:\*\?\"\/\<\>\|\0]`)
)

func NewFileService() FileService {
	return FileService(0)
}

func (s FileService) ValidName(fileName string) bool {
	return !strings.HasPrefix(fileName, " ") &&
		!strings.HasSuffix(fileName, " ") &&
		validFilenameRegexp.Match([]byte(fileName))
}

func (s FileService) PurgeFileName(fileName string) string {
	trimedFileName := strings.TrimSpace(fileName)

	return invalidFilenameCharacters.ReplaceAllString(trimedFileName, "_")
}

func (s FileService) CreateFilename(music *model.Music) string {
	nameBuilder := strings.Builder{}
	nameBuilder.WriteString(music.Author.Name)

	if len(music.Contributors) != 0 {
		nameBuilder.WriteString(" (")

		for i, contributors := range music.Contributors {
			nameBuilder.WriteString(contributors.Name)

			if len(music.Contributors)-1 != i {
				nameBuilder.WriteString(", ")
			}
		}

		nameBuilder.WriteString(")")
	}

	nameBuilder.WriteString(fmt.Sprintf(" - %s [%x]", music.Name, time.Now().UnixMicro()))

	return s.PurgeFileName(nameBuilder.String())
}

func (s FileService) MoveTo(from, to string) error {
	filename := pathModule.Base(from)
	newPath := pathModule.Join(to, filename)

	if f, err := os.Open(newPath); err == nil {
		f.Close()

		return fmt.Errorf("file (%s) already exists in (%s)", filename, to)
	}

	return os.Rename(from, newPath)
}

func (s FileService) IsExists(path string) bool {
	if f, err := os.Open(path); err == nil {
		f.Close()

		return true
	}

	return false
}
