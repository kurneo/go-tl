package postgres

import "time"

type Option func(postgres *Postgres)

func MaxPoolSize(size int) Option {
	return func(p *Postgres) {
		p.maxPoolSize = size
	}
}

func ConnAttempts(attempts int) Option {
	return func(p *Postgres) {
		p.connAttempts = attempts
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(p *Postgres) {
		p.connTimeout = timeout
	}
}

func SSLMode(enable bool) Option {
	return func(postgres *Postgres) {
		if enable {
			postgres.sslMode = "enable"
		} else {
			postgres.sslMode = "disable"
		}
	}
}

func Timezone(tz string) Option {
	return func(p *Postgres) {
		p.timezone = tz
	}
}
