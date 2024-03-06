package helper

import "strings"

type S3PathHelper struct {
}

func (s S3PathHelper) GetDirectoryPath(path string) string {
	p := strings.TrimRight(path, "\\/")
	if p == "" {
		return ""
	}
	return p
}

func NewS3PathHelper() S3PathHelper {
	return S3PathHelper{}
}
