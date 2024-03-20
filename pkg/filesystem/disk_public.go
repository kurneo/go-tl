package filesystem

import (
	"github.com/kurneo/go-template/pkg/filesystem/helper"
	"os"
	"strings"
)

type diskPublic struct {
	driver *driverLocal
	url    string
}

func (s diskPublic) FileExists(path string) (bool, error) {
	return s.driver.FileExists(path)
}

func (s diskPublic) DirExists(path string) (bool, error) {
	return s.driver.DirExists(path)
}

func (s diskPublic) Put(path string, content []byte) error {
	return s.driver.Put(path, content)
}

func (s diskPublic) Get(path string) ([]byte, error) {
	return s.driver.Get(path)
}

func (s diskPublic) MakeDir(path string, perm os.FileMode) error {
	return s.driver.MakeDir(path, perm)
}

func (s diskPublic) Delete(path string) error {
	return s.driver.Delete(path)
}

func (s diskPublic) Rename(from, to string) error {
	return s.driver.Rename(from, to)
}

func (s diskPublic) ListContents(path string) ([]File, []Directory, error) {
	return s.driver.ListContents(path)
}

func (s diskPublic) Move(from, to string) error {
	return s.driver.Move(from, to)
}

func (s diskPublic) Copy(from, to string) error {
	return s.driver.Copy(from, to)
}

func (s diskPublic) Mime(path string) string {
	return s.driver.Mime(path)
}

func (s diskPublic) RealPath(path string) string {
	return s.driver.RealPath(path)
}

func (s diskPublic) IsDir(path string) (bool, error) {
	return s.driver.IsDir(path)
}

func (s diskPublic) IsFile(path string) (bool, error) {
	return s.driver.IsFile(path)
}

func (s diskPublic) RealDirPath(path string) string {
	return s.driver.RealDirPath(path)
}

func (s diskPublic) Url(path string) string {
	if s.url == "" {
		return ""
	}
	return strings.TrimRight(s.url, "\\/") + "/" + path
}

func NewDiskPublic(prefix, separator, url string) DiskPublicContract {
	return &diskPublic{
		driver: &driverLocal{
			preFixer: helper.NewPreFixer(prefix, separator),
		},
		url: url,
	}
}
