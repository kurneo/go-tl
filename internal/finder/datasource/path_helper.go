package datasource

import (
	"reflect"
	"strings"
)

type PathHelper struct {
	separator string
}

func (o PathHelper) StripSlash(path string, vars ...interface{}) string {
	var cutSet = "\\/"
	if len(vars) > 0 && vars[0] != nil && reflect.ValueOf(vars[0]).Kind() == reflect.String {
		cutSet = vars[0].(string)
	}
	return strings.TrimRight(strings.TrimLeft(path, cutSet), cutSet)
}

func (o PathHelper) Concat(paths ...string) string {
	return strings.Join(paths, o.separator)
}

func (o PathHelper) DirPath(path string) string {
	split := strings.Split(path, o.separator)
	split = split[:len(split)-1]
	dirPath := o.Concat(split...)
	if dirPath == "." || dirPath == "" {
		return ""
	}
	return dirPath
}

func NewPathHelper(separator string) PathHelper {
	return PathHelper{separator: separator}
}
