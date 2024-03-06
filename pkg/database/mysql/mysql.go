package mysql

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kurneo/go-template/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

const (
	defaultMaxConnection = 1
	defaultConnAttempts  = 10
	defaultConnTimeout   = time.Second
)

type MySQL struct {
	maxConnection int
	connAttempts  int
	connTimeout   time.Duration
	charset       string
	parseTime     bool
	loc           string

	cfg           config.DBConn
	isTransaction bool
	db            *gorm.DB
	tx            *gorm.DB
}

func (m *MySQL) Close() error {
	if m.db != nil {
		db, err := m.db.DB()
		if err != nil {
			return err
		}
		err = db.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
func (m *MySQL) Connect() error {
	var err error
	for m.connAttempts > 0 {
		parseTime := "false"
		if m.parseTime {
			parseTime = "true"
		}
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
			m.cfg.Username,
			m.cfg.Password,
			m.cfg.Host,
			m.cfg.Port,
			m.cfg.DatabaseName,
			m.charset,
			parseTime,
			m.loc,
		)
		m.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			SkipDefaultTransaction: true,
			Logger:                 logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			break
		}

		log.Printf("trying to connect to Mysql, attempts left: %d", m.connAttempts)
		time.Sleep(m.connTimeout)
		m.connAttempts--
	}

	if err != nil {
		return err
	}
	return nil
}
func (m *MySQL) Begin() error {
	tx := m.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	m.isTransaction = true
	m.tx = tx
	return nil
}
func (m *MySQL) Commit() error {
	if m.tx == nil || !m.isTransaction {
		return errors.New("transaction is not start")
	}
	err := m.tx.Commit().Error
	if err != nil {
		return err
	}
	m.isTransaction = false
	m.tx = nil
	return nil
}
func (m *MySQL) Rollback() error {
	if m.tx == nil || !m.isTransaction {
		return errors.New("transaction is not start")
	}
	err := m.tx.Rollback().Error
	if err != nil {
		return err
	}
	m.isTransaction = false
	m.tx = nil
	return nil
}
func (m *MySQL) IsTransaction() bool {
	return m.isTransaction
}
func (m *MySQL) IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (m *MySQL) GetDB(ctx context.Context) *gorm.DB {
	if m.IsTransaction() {
		return m.tx.WithContext(ctx)
	}
	return m.db.WithContext(ctx)
}

func New(cfg config.DBConn, options ...Option) *MySQL {
	m := &MySQL{
		maxConnection: defaultMaxConnection,
		connAttempts:  defaultConnAttempts,
		connTimeout:   defaultConnTimeout,
		charset:       "utf8",
		parseTime:     true,
		loc:           "Local",
		cfg:           cfg,
	}
	for _, o := range options {
		o(m)
	}

	return m
}
