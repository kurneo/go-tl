package datasource

import (
	"errors"
	"github.com/kurneo/go-template/config"
	"github.com/kurneo/go-template/internal/finder/usecase"
	"github.com/kurneo/go-template/pkg/filesystem"
	"github.com/kurneo/go-template/pkg/filesystem/helper"
	"github.com/kurneo/go-template/pkg/logger"
)

func New(cfg config.Storage, l logger.Contract) (usecase.FinderRepositoryContract, error) {
	c, err := getConfig(cfg.Default, cfg)

	if err != nil {
		return nil, err
	}

	switch cfg.Default {
	case "public":
		return newLocal(c, l), nil
	default:
		return nil, errors.New("invalid disk: " + cfg.Default)
	}
}

func newLocal(c config.DiskCfg, l logger.Contract) usecase.FinderRepositoryContract {
	prefix := "files"
	separator := "/"

	pathHelper := NewPathHelper(separator)
	mainPreFixer := helper.NewPreFixer(
		pathHelper.Concat(pathHelper.StripSlash(c.Get("root")), prefix),
		separator,
	)
	return LocalRepo{
		mainStorage:    filesystem.NewPublicDisk(filesystem.NewLocalDriver(mainPreFixer, l), c.Get("url")),
		l:              l,
		pathHelper:     pathHelper,
		fileTypeHelper: NewFileTypeHelper(),
		Ignores:        []string{".gitkeep"},
	}
}

func getConfig(disk string, cfg config.Storage) (config.DiskCfg, error) {
	if c, ok := cfg.Disks[disk]; ok {
		return c, nil
	}
	return nil, errors.New("missing config for filesystem disk: " + disk)
}
