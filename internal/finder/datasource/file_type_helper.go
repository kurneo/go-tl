package datasource

import "strings"

type FileTypeHelper struct {
}

func (f FileTypeHelper) IsImageFile(mime string) bool {
	return strings.HasPrefix(mime, "image/")
}

func (f FileTypeHelper) IsVideoFile(mime string) bool {
	return strings.HasPrefix(mime, "video/")
}

func NewFileTypeHelper() FileTypeHelper {
	return FileTypeHelper{}
}
