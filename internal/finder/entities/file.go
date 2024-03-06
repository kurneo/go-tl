package entities

import (
	"time"
)

type File struct {
	Path      string
	Name      string
	ModTime   *time.Time
	Size      *int64
	Mime      *string
	Extension *string
	Url       string
}

func (file File) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type":      "file",
		"path":      file.Path,
		"name":      file.Name,
		"mod_time":  file.ModTime,
		"size":      file.Size,
		"mime":      file.Mime,
		"extension": file.Extension,
		"url":       file.Url,
	}
}
