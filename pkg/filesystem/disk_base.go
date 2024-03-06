package filesystem

import "os"

type disk struct {
	driver DriverContract
}

func (s disk) FileExists(path string) (bool, error) {
	return s.driver.FileExists(path)
}

func (s disk) DirExists(path string) (bool, error) {
	return s.driver.DirExists(path)
}

func (s disk) Put(path string, content []byte) error {
	return s.driver.Put(path, content)
}

func (s disk) Get(path string) ([]byte, error) {
	return s.driver.Get(path)
}

func (s disk) MakeDir(path string, perm os.FileMode) error {
	return s.driver.MakeDir(path, perm)
}

func (s disk) Delete(path string) error {
	return s.driver.Delete(path)
}

func (s disk) Rename(from, to string) error {
	return s.driver.Rename(from, to)
}

func (s disk) ListContents(path string) ([]File, []Directory, error) {
	return s.driver.ListContents(path)
}

func (s disk) Move(from, to string) error {
	return s.driver.Move(from, to)
}

func (s disk) Copy(from, to string) error {
	return s.driver.Copy(from, to)
}

func (s disk) Mime(path string) string {
	return s.driver.Mime(path)
}

func (s disk) RealPath(path string) string {
	return s.driver.RealPath(path)
}

func (s disk) IsDir(path string) (bool, error) {
	return s.driver.IsDir(path)
}

func (s disk) IsFile(path string) (bool, error) {
	return s.driver.IsFile(path)
}

func (s disk) RealDirPath(path string) string {
	return s.driver.RealDirPath(path)
}
