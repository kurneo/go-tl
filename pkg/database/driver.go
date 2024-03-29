package database

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

const (
	postgresDefaultMaxPoolSize  = 1
	postgresDefaultConnAttempts = 10
	postgresDefaultConnTimeout  = time.Second
)

type PgConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	DBName      string
	MaxPoolSize int
}

type Postgres struct {
	c            PgConfig
	connAttempts int
	connTimeout  time.Duration

	isTransaction bool
	db            *gorm.DB
	tx            *gorm.DB
}

func (p *Postgres) Close() error {
	if p.db != nil {
		db, err := p.db.DB()
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

func (p *Postgres) Connect() error {
	var err error

	//maxPoolSize := viper.GetInt("POSTGRES_DB_MAX_POOL_SIZE")
	if p.c.MaxPoolSize == 0 {
		p.c.MaxPoolSize = postgresDefaultMaxPoolSize
	}

	for p.connAttempts > 0 {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d",
			p.c.Host,
			p.c.User,
			p.c.Password,
			p.c.DBName,
			p.c.Port,
		)

		p.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			SkipDefaultTransaction: true,
			Logger:                 logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			break
		}

		log.Printf("trying to connect to Postgres, attempts left: %d", p.connAttempts)
		time.Sleep(p.connTimeout)
		p.connAttempts--
	}

	if err != nil {
		return err
	}

	db, err := p.db.DB()

	if err != nil {
		return err
	}

	db.SetMaxOpenConns(p.c.MaxPoolSize)

	return nil
}

func (p *Postgres) IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (p *Postgres) Begin() error {
	tx := p.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	p.isTransaction = true
	p.tx = tx
	return nil
}

func (p *Postgres) Commit() error {
	if p.tx == nil || !p.isTransaction {
		return errors.New("transaction is not start")
	}
	err := p.tx.Commit().Error
	if err != nil {
		return err
	}
	p.isTransaction = false
	p.tx = nil
	return nil
}

func (p *Postgres) Rollback() error {
	if p.tx == nil || !p.isTransaction {
		return errors.New("transaction is not start")
	}
	err := p.tx.Rollback().Error
	if err != nil {
		return err
	}
	p.isTransaction = false
	p.tx = nil
	return nil
}

func (p *Postgres) IsTransaction() bool {
	return p.isTransaction
}

func (p *Postgres) GetDB(ctx context.Context) *gorm.DB {
	if p.IsTransaction() {
		return p.tx.WithContext(ctx)
	}
	return p.db.WithContext(ctx)
}

func newPostgres(c PgConfig) *Postgres {
	return &Postgres{
		c:            c,
		connAttempts: postgresDefaultConnAttempts,
		connTimeout:  postgresDefaultConnTimeout,
	}
}

const (
	mySqlDefaultMaxConnection = 1
	mySqlDefaultConnAttempts  = 10
	mySqlDefaultConnTimeout   = time.Second
)

type MySqlConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type MySQL struct {
	c            MySqlConfig
	connAttempts int
	connTimeout  time.Duration

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
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s",
			m.c.User,
			m.c.Password,
			m.c.Host,
			m.c.Port,
			m.c.DBName,
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

func newMySql(c MySqlConfig) *MySQL {
	return &MySQL{
		c:            c,
		connAttempts: mySqlDefaultConnAttempts,
		connTimeout:  mySqlDefaultConnTimeout,
	}
}
