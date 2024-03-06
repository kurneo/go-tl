package datasource

import (
	"github.com/kurneo/go-template/pkg/filesystem"
	"github.com/kurneo/go-template/pkg/logger"
)

type S3Repo struct {
	mainStorage filesystem.S3DiskContract
	l           logger.Contract
	pathHelper  PathHelper
	Ignores     []string
}
