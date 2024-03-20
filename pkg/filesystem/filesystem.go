package filesystem

import (
	"os"
	"time"
)

const (
	DriverLocal = "local"
	DriverS3    = "s3"
	DriverMinio = "minio"
)

type DriverContract interface {
	FileExists(path string) (bool, error)
	DirExists(path string) (bool, error)
	Put(path string, content []byte) error
	Get(path string) ([]byte, error)
	MakeDir(path string, perm os.FileMode) error
	Delete(path string) error
	Rename(from, to string) error
	ListContents(path string) ([]File, []Directory, error)
	Move(from, to string) error
	Copy(from, to string) error
	Mime(path string) string
	RealPath(path string) string
	RealDirPath(path string) string
	IsDir(path string) (bool, error)
	IsFile(path string) (bool, error)
}

type DiskLocalContract DriverContract

type DiskPublicContract interface {
	DriverContract
	Url(path string) string
}

type DiskS3Contract interface {
	DriverContract
	Url(path string) string
}

type File struct {
	Path      string
	Name      string
	ModTime   *time.Time
	Size      *int64
	Mime      *string
	Extension *string
}

type Directory struct {
	Path    string
	Name    string
	ModTime *time.Time
}
