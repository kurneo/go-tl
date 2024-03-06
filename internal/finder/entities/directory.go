package entities

import "time"

type Directory struct {
	Path    string     `json:"path"`
	Name    string     `json:"name"`
	ModTime *time.Time `json:"mod_time"`
}

func (dir Directory) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type":     "dir",
		"path":     dir.Path,
		"name":     dir.Name,
		"mod_time": dir.ModTime,
	}
}
