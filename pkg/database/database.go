package database

import (
	"context"
	"errors"
	"github.com/kurneo/go-template/config"
	"github.com/kurneo/go-template/pkg/database/mysql"
	"github.com/kurneo/go-template/pkg/database/postgres"
	"gorm.io/gorm"
	"sync"
)

type Contract interface {
	Close() error
	Connect() error
	Begin() error
	Commit() error
	Rollback() error
	IsTransaction() bool
	IsNotFound(err error) bool
	GetDB(ctx context.Context) *gorm.DB
}

var (
	DriverPostgres = "postgres"
	DriverMysql    = "mysql"
	dbInstance     Contract
	dbOnce         sync.Once
)

func New(cfg config.DB, options ...interface{}) (Contract, error) {
	var err error = nil
	c, ok := cfg.Connections[cfg.Default]

	if !ok {
		return nil, errors.New("cannot find config for: " + cfg.Default)
	}

	dbOnce.Do(func() {
		switch cfg.Default {
		case DriverPostgres:
			opts := make([]postgres.Option, 0)
			for _, o := range options {
				opts = append(opts, o.(postgres.Option))
			}
			dbInstance = postgres.New(c, opts...)
			break
		case DriverMysql:
			opts := make([]mysql.Option, 0)
			for _, o := range options {
				opts = append(opts, o.(mysql.Option))
			}
			dbInstance = mysql.New(c, opts...)
		default:
			dbInstance = nil
		}

		if dbInstance == nil {
			err = errors.New("invalid database connection: " + cfg.Default)
		}

		if errConnect := dbInstance.Connect(); errConnect != nil {
			err = errConnect
		}
	})

	return dbInstance, err
}
