package usecase

import (
	"github.com/kurneo/go-template/internal/finder/entities"
	errPkg "github.com/kurneo/go-template/pkg/error"
	"mime/multipart"
)

type (
	FinderUseCaseContract interface {
		GetContents(path, sortField, sortDir string) ([]entities.File, []entities.Directory, errPkg.Contract)
		CreateDirectory(dto CreateDirDTO) errPkg.Contract
		Thumbnail(path string) ([]byte, errPkg.Contract)
		Preview(path string) ([]byte, string, errPkg.Contract)
		Upload(path string, file *multipart.FileHeader) errPkg.Contract
		Rename(dto RenameDTO) errPkg.Contract
		Delete(dto DeleteDTO) errPkg.Contract
		Copy(dto CopyDTO) errPkg.Contract
		Cut(dto CutDTO) errPkg.Contract
	}

	FinderRepositoryContract interface {
		GetContents(path, sortField, sortDir string) ([]entities.File, []entities.Directory, errPkg.Contract)
		CreateDirectory(path string) errPkg.Contract
		FileExist(path string) (bool, errPkg.Contract)
		DirExist(path string) (bool, errPkg.Contract)
		Thumbnail(path string) ([]byte, errPkg.Contract)
		Preview(path string) ([]byte, string, errPkg.Contract)
		Upload(path string, file *multipart.FileHeader) (string, errPkg.Contract)
		Rename(items map[string]string) errPkg.Contract
		Delete(items []string) errPkg.Contract
		Copy(items map[string]string) errPkg.Contract
		Cut(items map[string]string) errPkg.Contract
	}

	CreateDirDTO interface {
		GetName() string
		GetPath() string
	}

	RenameDTO interface {
		GetItems() map[string]string
	}

	DeleteDTO interface {
		GetItems() []string
	}

	CopyDTO interface {
		GetItems() []string
		GetPath() string
	}

	CutDTO interface {
		CopyDTO
	}
)
