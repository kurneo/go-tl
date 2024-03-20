package filesystem

import (
	"github.com/kurneo/go-template/pkg/filesystem/helper"
	"os"
)

type diskLocal struct {
	driver *driverLocal
}

func (s diskLocal) FileExists(path string) (bool, error) {
	return s.driver.FileExists(path)
}

func (s diskLocal) DirExists(path string) (bool, error) {
	return s.driver.DirExists(path)
}

func (s diskLocal) Put(path string, content []byte) error {
	return s.driver.Put(path, content)
}

func (s diskLocal) Get(path string) ([]byte, error) {
	return s.driver.Get(path)
}

func (s diskLocal) MakeDir(path string, perm os.FileMode) error {
	return s.driver.MakeDir(path, perm)
}

func (s diskLocal) Delete(path string) error {
	return s.driver.Delete(path)
}

func (s diskLocal) Rename(from, to string) error {
	return s.driver.Rename(from, to)
}

func (s diskLocal) ListContents(path string) ([]File, []Directory, error) {
	return s.driver.ListContents(path)
}

func (s diskLocal) Move(from, to string) error {
	return s.driver.Move(from, to)
}

func (s diskLocal) Copy(from, to string) error {
	return s.driver.Copy(from, to)
}

func (s diskLocal) Mime(path string) string {
	return s.driver.Mime(path)
}

func (s diskLocal) RealPath(path string) string {
	return s.driver.RealPath(path)
}

func (s diskLocal) IsDir(path string) (bool, error) {
	return s.driver.IsDir(path)
}

func (s diskLocal) IsFile(path string) (bool, error) {
	return s.driver.IsFile(path)
}

func (s diskLocal) RealDirPath(path string) string {
	return s.driver.RealDirPath(path)
}

func NewDiskLocal(prefix, separator string) DiskLocalContract {
	return &diskLocal{
		driver: &driverLocal{
			preFixer: helper.NewPreFixer(prefix, separator),
		},
	}
}
