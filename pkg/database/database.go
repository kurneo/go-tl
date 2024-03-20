package database

import (
	"context"
	"errors"
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

type Config struct {
	Driver string
	PgSql  PgConfig
	MySql  MySqlConfig
}

const (
	DriverPostgres = "pgsql"
	DriverMysql    = "mysql"
)

var (
	dbInstance Contract
	dbOnce     sync.Once
)

func New(c Config) (Contract, error) {
	var err error = nil

	if c.Driver == "" || (c.Driver != DriverMysql && c.Driver != DriverPostgres) {
		return nil, errors.New("default database driver is invalid")
	}

	dbOnce.Do(func() {
		switch c.Driver {
		case DriverPostgres:
			dbInstance = newPostgres(c.PgSql)
			break
		case DriverMysql:
			dbInstance = newMySql(c.MySql)
		}

		if errConnect := dbInstance.Connect(); errConnect != nil {
			err = errConnect
		}
	})

	return dbInstance, err
}
