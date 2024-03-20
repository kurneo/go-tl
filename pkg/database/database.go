package database

import (
	"context"
	"errors"
	"github.com/spf13/viper"
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

const (
	DriverPostgres = "pgsql"
	DriverMysql    = "mysql"
)

var (
	dbInstance Contract
	dbOnce     sync.Once
)

func New() (Contract, error) {
	var err error = nil

	d := viper.GetString("DB_DRIVER")
	if d == "" || (d != DriverMysql && d != DriverPostgres) {
		return nil, errors.New("default database driver is invalid")
	}

	dbOnce.Do(func() {
		switch d {
		case DriverPostgres:
			dbInstance = newPostgres()
			break
		case DriverMysql:
			dbInstance = newMySql()
		}

		if errConnect := dbInstance.Connect(); errConnect != nil {
			err = errConnect
		}
	})

	return dbInstance, err
}
