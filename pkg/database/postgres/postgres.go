package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/kurneo/go-template/config"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

const (
	defaultMaxPoolSize  = 1
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration
	sslMode      string
	timezone     string

	cfg           config.DBConn
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
	for p.connAttempts > 0 {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			p.cfg.Host,
			p.cfg.Username,
			p.cfg.Password,
			p.cfg.DatabaseName,
			p.cfg.Port,
			p.sslMode,
			p.timezone,
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

	db.SetMaxOpenConns(p.maxPoolSize)

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

func New(cfg config.DBConn, options ...Option) *Postgres {
	pg := &Postgres{
		maxPoolSize:  defaultMaxPoolSize,
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
		cfg:          cfg,
		sslMode:      "disable",
		timezone:     "UTC",
	}

	for _, option := range options {
		option(pg)
	}

	return pg
}
