package file

import (
	"regexp"
	"strings"
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
