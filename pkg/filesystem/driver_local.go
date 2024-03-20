package filesystem

import (
	"fmt"
	"github.com/kurneo/go-template/pkg/filesystem/helper"
	cp "github.com/otiai10/copy"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

type driverLocal struct {
	preFixer helper.PathPreFixer
}

func (d driverLocal) FileExists(path string) (bool, error) {
	stat, err := os.Stat(d.preFixer.PrefixPath(path))

	if stat != nil {
		return stat.IsDir() == false, nil
	}

	if !os.IsNotExist(err) {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (d driverLocal) DirExists(path string) (bool, error) {
	stat, err := os.Stat(d.preFixer.PrefixPath(path))

	if stat != nil {
		return stat.IsDir(), nil
	}

	if !os.IsNotExist(err) {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (d driverLocal) Put(path string, content []byte) error {
	f, err := os.Create(d.preFixer.PrefixPath(path))
	if err != nil {
		return err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	_, err = f.WriteString(string(content))

	if err != nil {
		return err
	}

	return nil
}

func (d driverLocal) Get(path string) ([]byte, error) {
	file, err := os.Open(d.preFixer.PrefixPath(path))
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (d driverLocal) MakeDir(path string, perm os.FileMode) error {
	err := os.MkdirAll(d.preFixer.PrefixPath(path), perm)
	if err != nil {
		return err
	}
	return nil
}

func (d driverLocal) Delete(path string) error {
	err := os.RemoveAll(d.preFixer.PrefixPath(path))
	if err != nil {
		return err
	}
	return nil
}

func (d driverLocal) Rename(from, to string) error {
	err := os.Rename(d.preFixer.PrefixPath(from), d.preFixer.PrefixPath(to))
	if err != nil {
		return err
	}
	return nil
}

func (d driverLocal) ListContents(path string) ([]File, []Directory, error) {
	reader, err := os.Open(d.preFixer.PrefixPath(path))
	if err != nil {
		return nil, nil, err
	}
	items, err := reader.Readdir(0)

	if err != nil {
		return nil, nil, err
	}

	files := make([]File, 0)
	directories := make([]Directory, 0)

	for _, item := range items {
		itemP := strings.TrimLeft(path+"/"+item.Name(), "\\/")
		if item.IsDir() {
			t := item.ModTime()
			directories = append(directories, Directory{
				Path:    itemP,
				ModTime: &t,
				Name:    item.Name(),
			})
		} else {
			t := item.ModTime()
			s := item.Size()
			e := filepath.Ext(item.Name())
			file := File{
				Path:      itemP,
				ModTime:   &t,
				Name:      item.Name(),
				Size:      &s,
				Extension: &e,
			}
			m := d.Mime(itemP)
			file.Mime = &m
			files = append(files, file)
		}
	}

	return files, directories, nil
}

func (d driverLocal) Move(from, to string) error {
	return d.Rename(from, to)
}

func (d driverLocal) Copy(from, to string) error {
	err := cp.Copy(d.preFixer.PrefixPath(from), d.preFixer.PrefixPath(to))
	if err != nil {
		return err
	}
	return nil
}

func (d driverLocal) Mime(path string) string {
	return mime.TypeByExtension(filepath.Ext(path))
}

func (d driverLocal) RealPath(path string) string {
	return d.preFixer.PrefixPath(path)
}

func (d driverLocal) RealDirPath(path string) string {
	return d.preFixer.PrefixDirectoryPath(path)
}

func (d driverLocal) IsDir(path string) (bool, error) {
	fi, err := os.Stat(d.preFixer.PrefixPath(path))
	if err != nil {
		return false, err
	} else if fi.IsDir() {
		return true, nil
	}
	return false, nil
}

func (d driverLocal) IsFile(path string) (bool, error) {
	fi, err := os.Stat(d.preFixer.PrefixPath(path))
	if err != nil {
		return false, err
	} else if !fi.IsDir() {
		return true, nil
	}
	return false, nil
}

func NewDriverLocal(prefix, separator string) DriverContract {
	return driverLocal{
		preFixer: helper.NewPreFixer(prefix, separator),
	}
}
