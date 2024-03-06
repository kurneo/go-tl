package mysql

import "time"

type Option func(mysql *MySQL)

func MaxConnection(size int) Option {
	return func(p *MySQL) {
		p.maxConnection = size
	}
}

func ConnAttempts(attempts int) Option {
	return func(p *MySQL) {
		p.connAttempts = attempts
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(p *MySQL) {
		p.connTimeout = timeout
	}
}

func Charset(s string) Option {
	return func(mysql *MySQL) {
		mysql.charset = s
	}
}

func ParseTime(v bool) Option {
	return func(p *MySQL) {
		p.parseTime = v
	}
}
func Loc(v string) Option {
	return func(p *MySQL) {
		p.loc = v
	}
}
