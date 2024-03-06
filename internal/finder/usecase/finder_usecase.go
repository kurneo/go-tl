package usecase

import (
	"errors"
	"github.com/kurneo/go-template/internal/finder/entities"
	"github.com/kurneo/go-template/pkg/error"
	"github.com/kurneo/go-template/pkg/logger"
	"mime/multipart"
	path2 "path"
	"strings"
)

var (
	errPathNotExist     = errors.New("path not exist")
	errPathAlreadyExist = errors.New("path already exist")
)

type FinderUseCase struct {
	log logger.Contract
	r   FinderRepositoryContract
}

func (u FinderUseCase) GetContents(path, sortField, sortDir string) ([]entities.File, []entities.Directory, error.Contract) {
	exist, err := u.r.DirExist(path)
	if err != nil {
		return nil, nil, err
	}
	if !exist {
		return nil, nil, error.NewDomain(errPathNotExist)
	}
	return u.r.GetContents(path, sortField, sortDir)
}

func (u FinderUseCase) CreateDirectory(dto CreateDirDTO) error.Contract {
	path := dto.GetName()
	if u.trimSlash(dto.GetPath()) != "" {
		path = u.concatPaths(u.trimSlash(dto.GetPath()), path)
	}
	exist, err := u.r.DirExist(path)
	if err != nil {
		return err
	}
	if exist {
		return error.NewDomain(errPathAlreadyExist)
	}
	err = u.r.CreateDirectory(path)
	if err != nil {
		return err
	}
	return nil
}

func (u FinderUseCase) Thumbnail(path string) ([]byte, error.Contract) {
	exist, err := u.r.FileExist(path)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, error.NewDomain(errPathNotExist)
	}
	return u.r.Thumbnail(path)
}

func (u FinderUseCase) Preview(path string) ([]byte, string, error.Contract) {
	exist, err := u.r.FileExist(path)
	if err != nil {
		return nil, "", err
	}
	if !exist {
		return nil, "", error.NewDomain(errPathNotExist)
	}

	return u.r.Preview(path)
}

func (u FinderUseCase) Upload(path string, file *multipart.FileHeader) error.Contract {
	filePath := file.Filename
	if u.trimSlash(path) != "" {
		filePath = u.concatPaths(u.trimSlash(path), filePath)
	}

	exist, err := u.r.FileExist(filePath)
	if err != nil {
		return err
	}

	if exist {
		return error.NewDomain(errors.New("path exists"))
	}

	_, err = u.r.Upload(path, file)

	if err != nil {
		return err
	}

	return nil
}

func (u FinderUseCase) Rename(dto RenameDTO) error.Contract {
	items := make(map[string]string)

	for oldPath, newName := range dto.GetItems() {
		paths := []string{newName}
		dir := u.dirPath(u.trimSlash(oldPath))
		if dir != "" {
			paths = append([]string{dir}, paths...)
		}
		newPath := u.concatPaths(paths...)

		if oldPath == newPath {
			continue
		}

		exist, err := u.r.FileExist(oldPath)
		if err != nil {
			return err
		}
		if !exist {
			return error.NewDomain(errPathNotExist)
		}

		exist, err = u.r.FileExist(newPath)
		if err != nil {
			return err
		}
		if exist {
			return error.NewDomain(errPathAlreadyExist)
		}

		items[oldPath] = newPath
	}

	if err := u.r.Rename(items); err != nil {
		return err
	}

	return nil
}

func (u FinderUseCase) Delete(dto DeleteDTO) error.Contract {
	items := make([]string, 0)
	for _, path := range dto.GetItems() {
		exist, err := u.r.FileExist(path)
		if err != nil {
			return err
		}
		if !exist {
			return error.NewDomain(errPathNotExist)
		}

		items = append(items, path)
	}

	err := u.r.Delete(items)
	if err != nil {
		return err
	}

	return nil
}

func (u FinderUseCase) Copy(dto CopyDTO) error.Contract {
	items := make(map[string]string)
	for _, path := range dto.GetItems() {
		exist, err := u.r.FileExist(path)
		if err != nil {
			return err
		}
		if !exist {
			return error.NewDomain(errPathNotExist)
		}

		newPath := u.concatPaths(dto.GetPath(), u.baseName(path))

		exist, err = u.r.FileExist(newPath)
		if err != nil {
			return err
		}
		if exist {
			return error.NewDomain(errPathAlreadyExist)
		}
		items[path] = newPath
	}

	err := u.r.Copy(items)
	if err != nil {
		return err
	}
	return nil
}

func (u FinderUseCase) Cut(dto CutDTO) error.Contract {
	items := make(map[string]string)
	for _, path := range dto.GetItems() {
		exist, err := u.r.FileExist(path)
		if err != nil {
			return err
		}
		if !exist {
			return error.NewDomain(errPathNotExist)
		}

		newPath := u.concatPaths(dto.GetPath(), u.baseName(path))

		exist, err = u.r.FileExist(newPath)
		if err != nil {
			return err
		}
		if exist {
			return error.NewDomain(errPathAlreadyExist)
		}
		items[path] = newPath
	}

	err := u.r.Cut(items)
	if err != nil {
		return err
	}

	return nil
}

func (u FinderUseCase) trimSlash(p string) string {
	cutSet := "/"
	return strings.TrimLeft(strings.TrimRight(p, cutSet), cutSet)
}

func (u FinderUseCase) dirPath(path string) string {
	split := strings.Split(path, "/")
	split = split[:len(split)-1]
	dirPath := strings.Join(split, "/")
	if dirPath == "." || dirPath == "" {
		return ""
	}
	return dirPath
}

func (u FinderUseCase) concatPaths(paths ...string) string {
	return strings.Join(paths, "/")
}

func (u FinderUseCase) baseName(path string) string {
	return path2.Base(path)
}

func New(l logger.Contract, r FinderRepositoryContract) FinderUseCaseContract {
	return FinderUseCase{
		r:   r,
		log: l,
	}
}
