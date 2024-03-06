package datasource

import (
	"bytes"
	"github.com/disintegration/imaging"
	"github.com/kurneo/go-template/internal/finder/entities"
	pkgErr "github.com/kurneo/go-template/pkg/error"
	"github.com/kurneo/go-template/pkg/filesystem"
	"github.com/kurneo/go-template/pkg/logger"
	"github.com/kurneo/go-template/pkg/support/slices"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalRepo struct {
	mainStorage    filesystem.PublicDiskContract
	l              logger.Contract
	pathHelper     PathHelper
	fileTypeHelper FileTypeHelper
	Ignores        []string
}

func (r LocalRepo) GetContents(path, sortField, sortDir string) ([]entities.File, []entities.Directory, pkgErr.Contract) {
	files, dirs, err := r.mainStorage.ListContents(path)
	if err != nil {
		return nil, nil, pkgErr.NewDatasource(err)
	}
	return slices.Map[filesystem.File, entities.File](
			slices.Filter[filesystem.File](files, func(file filesystem.File) bool {
				return !r.isIgnore(file.Name)
			}),
			func(file filesystem.File) entities.File {
				return entities.File{
					Path:      file.Path,
					Name:      file.Name,
					ModTime:   file.ModTime,
					Size:      file.Size,
					Mime:      file.Mime,
					Extension: file.Extension,
					Url:       r.mainStorage.Url(file.Path),
				}
			},
		), slices.Map[filesystem.Directory, entities.Directory](
			slices.Filter[filesystem.Directory](dirs, func(file filesystem.Directory) bool {
				return !r.isIgnore(file.Name)
			}),
			func(directory filesystem.Directory) entities.Directory {
				return entities.Directory{
					Path:    directory.Path,
					Name:    directory.Name,
					ModTime: directory.ModTime,
				}
			},
		), nil
}

func (r LocalRepo) FileExist(path string) (bool, pkgErr.Contract) {
	exist, err := r.mainStorage.FileExists(path)
	if err != nil {
		return false, pkgErr.NewDatasource(err)
	}
	return exist, nil
}

func (r LocalRepo) DirExist(path string) (bool, pkgErr.Contract) {
	exist, err := r.mainStorage.DirExists(path)
	if err != nil {
		return false, pkgErr.NewDatasource(err)
	}
	return exist, nil
}

func (r LocalRepo) CreateDirectory(path string) pkgErr.Contract {
	err := r.mainStorage.MakeDir(path, 0777)
	if err != nil {
		return pkgErr.NewDatasource(err)
	}
	return nil
}

func (r LocalRepo) Thumbnail(path string) ([]byte, pkgErr.Contract) {
	src, err := imaging.Open(r.mainStorage.RealPath(path))
	if err != nil {
		return nil, pkgErr.NewDatasource(err)
	}
	var buf bytes.Buffer
	src = imaging.Thumbnail(src, 100, 100, imaging.CatmullRom)
	err = imaging.Encode(&buf, src, imaging.JPEG)
	if err != nil {
		return nil, pkgErr.NewDatasource(err)
	}
	return buf.Bytes(), nil
}

func (r LocalRepo) Preview(path string) ([]byte, string, pkgErr.Contract) {
	file, err := r.mainStorage.Get(path)

	if err != nil {
		return nil, "", pkgErr.NewDatasource(err)
	}

	mime := r.mainStorage.Mime(path)

	if err != nil {
		return nil, "", pkgErr.NewDatasource(err)
	}

	return file, mime, nil
}

func (r LocalRepo) Upload(path string, file *multipart.FileHeader) (string, pkgErr.Contract) {
	src, err := file.Open()
	if err != nil {
		r.l.Error(err)
		return "", pkgErr.NewDatasource(err)
	}
	defer func() {
		err := src.Close()
		if err != nil {
			r.l.Error(err)
		}
	}()

	// Destination
	filePath := r.pathHelper.Concat(r.pathHelper.StripSlash(path), file.Filename)

	dst, err := os.Create(r.mainStorage.RealPath(filePath))
	if err != nil {
		r.l.Error(err)
		return "", pkgErr.NewDatasource(err)
	}
	defer func() {
		err := dst.Close()
		if err != nil {
			r.l.Error(err)
		}
	}()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		r.l.Error(err)
		return "", pkgErr.NewDatasource(err)
	}

	return filePath, nil
}

func (r LocalRepo) IsDir(path string) (bool, pkgErr.Contract) {
	isDir, err := r.mainStorage.IsDir(path)
	if err != nil {
		return false, pkgErr.NewDatasource(err)
	}
	return isDir, nil
}

func (r LocalRepo) Rename(items map[string]string) pkgErr.Contract {
	for oldPath, newPath := range items {
		err := r.mainStorage.Rename(oldPath, newPath)
		if err != nil {
			return pkgErr.NewDatasource(err)
		}
	}
	return nil
}

func (r LocalRepo) Delete(items []string) pkgErr.Contract {
	for _, path := range items {
		err := r.mainStorage.Delete(path)
		if err != nil {
			return pkgErr.NewDatasource(err)
		}
	}
	return nil
}

func (r LocalRepo) Copy(items map[string]string) pkgErr.Contract {
	for oldPath, newPath := range items {
		err := r.mainStorage.Copy(oldPath, newPath)
		if err != nil {
			return pkgErr.NewDatasource(err)
		}
	}
	return nil
}

func (r LocalRepo) Cut(items map[string]string) pkgErr.Contract {
	for oldPath, newPath := range items {
		err := r.mainStorage.Move(oldPath, newPath)
		if err != nil {
			return pkgErr.NewDatasource(err)
		}
	}
	return nil
}

func (r LocalRepo) BaseName(path string) string {
	return filepath.Base(path)
}

func (r LocalRepo) DirPath(path string) string {
	return r.pathHelper.DirPath(path)
}

func (r LocalRepo) StripSlash(path string) string {
	return r.pathHelper.StripSlash(path)
}

func (r LocalRepo) ConcatPaths(paths ...string) string {
	return r.pathHelper.Concat(paths...)
}

func (r LocalRepo) isIgnore(name string) bool {
	for _, v := range r.Ignores {
		if v == name {
			return true
		}
	}
	return false
}
