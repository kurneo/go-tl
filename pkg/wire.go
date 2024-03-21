package pkg

import (
	"github.com/google/wire"
	"github.com/kurneo/go-template/pkg/cache"
	"github.com/kurneo/go-template/pkg/database"
	"github.com/kurneo/go-template/pkg/hashing"
	"github.com/kurneo/go-template/pkg/jwt"
	logPkg "github.com/kurneo/go-template/pkg/log"
	"github.com/kurneo/go-template/pkg/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"log"
	"strings"
	"time"
)

// WireSet DI export from pkg
var WireSet = wire.NewSet(
	ResolveCacheInstance,
	ResolveDatabaseInstance,
	ResolveLogInstance,
	ResolveTokenManager,
	ResolveJWTMiddlewareFunc,
	ResolveHashingInstance,
	ResolveEcho,
)

// ResolveCacheInstance resolve dependencies and create cache instance
func ResolveCacheInstance() cache.Contact {
	cfg := cache.Config{
		Driver: viper.GetString("CACHE_DRIVER"),
		Redis: struct {
			Host string
			Port int
			DB   int
		}{
			Host: viper.GetString("CACHE_REDIS_HOST"),
			Port: viper.GetInt("CACHE_REDIS_PORT"),
			DB:   viper.GetInt("CACHE_REDIS_DB"),
		},
		InMemory: struct {
			DefaultExpiration      time.Duration
			DefaultCleanUpInterval time.Duration
		}{
			DefaultExpiration:      viper.GetDuration("CACHE_IN_MEMORY_DEFAULT_EXPIRATION"),
			DefaultCleanUpInterval: viper.GetDuration("CACHE_IN_MEMORY_CLEANUP_INTERVAL"),
		},
	}
	c, err := cache.New(cfg)
	if err != nil {
		log.Fatalf("init cache error: %s", err)
	}
	return c
}

// ResolveDatabaseInstance resolve global database instance
func ResolveDatabaseInstance() database.Contract {
	c := database.Config{
		Driver: viper.GetString("DB_DRIVER"),
		PgSql: struct {
			Host        string
			Port        int
			User        string
			Password    string
			DBName      string
			MaxPoolSize int
		}{
			Host:     viper.GetString("POSTGRES_DB_HOST"),
			Port:     viper.GetInt("POSTGRES_DB_PORT"),
			User:     viper.GetString("POSTGRES_DB_USER"),
			Password: viper.GetString("POSTGRES_DB_PASSWORD"),
			DBName:   viper.GetString("POSTGRES_DB_NAME"),
		},
		MySql: struct {
			Host     string
			Port     int
			User     string
			Password string
			DBName   string
		}{
			Host:     viper.GetString("MYSQL_DB_HOST"),
			Port:     viper.GetInt("MYSQL_DB_PORT"),
			User:     viper.GetString("MYSQL_DB_USER"),
			Password: viper.GetString("MYSQL_DB_PASSWORD"),
			DBName:   viper.GetString("MYSQL_DB_NAME"),
		},
	}

	d, err := database.New(c)
	if err != nil {
		log.Fatalf("init database error: %s", err)
	}
	return d
}

// ResolveLogInstance resolve global log instance
func ResolveLogInstance() logPkg.Contract {
	c := logPkg.Config{
		Channel: viper.GetString("LOG_DEFAULT_CHANNEL"),
		Daily: struct {
			FileName string
			Level    string
		}{
			FileName: viper.GetString("LOG_DAILY_FILE_NAME"),
			Level:    viper.GetString("LOG_DAILY_LOG_LEVEL"),
		},
		Singe: struct {
			FileName string
			Level    string
		}{
			FileName: viper.GetString("LOG_SINGLE_FILE_NAME"),
			Level:    viper.GetString("LOG_SINGLE_LOG_LEVEL"),
		},
		StdOut: struct{ Level string }{
			Level: viper.GetString("LOG_STDOUT_LOG_LEVEL"),
		},
		TeleHookConfig: struct {
			Enable   bool
			BotToken string
			ChatID   string
			Level    string
			Mentions string
		}{
			Enable:   viper.GetBool("LOG_HOOK_TELE_ENABLE"),
			BotToken: viper.GetString("LOG_HOOK_TELE_BOT_TOKEN"),
			ChatID:   viper.GetString("LOG_HOOK_TELE_CHAT_ID"),
			Level:    viper.GetString("LOG_HOOK_TELE_LEVEL"),
			Mentions: viper.GetString("LOG_HOOK_TELE_MENTIONS"),
		},
	}
	l, err := logPkg.New(c)
	if err != nil {
		log.Fatalf("init logger error: %s", err)
	}
	return l
}

// ResolveTokenManager resolve global jwt token manager instance
func ResolveTokenManager(c cache.Contact) *jwt.TokenManager[int64] {
	cfg := jwt.Config{
		Secret:  viper.GetString("JWT_SECRET"),
		Timeout: viper.GetInt("JWT_TOKEN_TIMEOUT"),
	}
	return jwt.NewTokenManager[int64](c, cfg)
}

// ResolveJWTMiddlewareFunc resolve global echo jwt middleware
func ResolveJWTMiddlewareFunc(t *jwt.TokenManager[int64]) echo.MiddlewareFunc {
	return middlewares.JwtMiddleware(t)
}

// ResolveHashingInstance  resolve global hashing instance
func ResolveHashingInstance() hashing.Contact {
	c := hashing.Config{
		Driver: viper.GetString("HASHING_DRIVER"),
		Bcrypt: struct{ Cost int }{Cost: viper.GetInt("HASHING_BCRYPT_COST")},
		Argon2: struct {
			Memory      uint32
			Iterations  uint32
			Parallelism uint8
			SaltLength  uint32
			KeyLength   uint32
		}{
			Memory:      viper.GetUint32("HASHING_ARGON2_MEMORY"),
			Iterations:  viper.GetUint32("HASHING_ARGON2_ITERATIONS"),
			Parallelism: uint8(viper.GetInt("HASHING_ARGON2_PARALLELISM")),
			SaltLength:  viper.GetUint32("HASHING_ARGON2_SALT_LENGTH"),
			KeyLength:   viper.GetUint32("HASHING_ARGON2_KEY_LENGTH"),
		},
	}

	s, err := hashing.New(c)
	if err != nil {
		log.Fatalf("init hashing error: %s", err)
	}
	return s
}

// ResolveEcho resolve global echo instance
func ResolveEcho(jwtMiddleware echo.MiddlewareFunc) *echo.Echo {
	echoApp := echo.New()
	// configure global middleware here

	c := viper.GetString("HTTP_HEADER_ALLOW_ORIGIN")
	l := viper.GetInt("HTTP_HEADER_GZIP_LEVEL")
	e := viper.GetString("HTTP_HEADER_EXPOSE")
	r := viper.GetFloat64("HTTP_THROTTLE_RATE_LIMIT")
	d := viper.GetDuration("HTTP_THROTTLE_DECAY")

	if r == 0 {
		r = 300
	}

	if d < time.Minute {
		d = time.Minute
	}

	echoApp.Use(
		middlewares.CorsMiddleware(strings.Split(c, ",")),
		middlewares.RateLimiterMiddleware(r, d),
		middlewares.GzipMiddleware(l),
		middlewares.AddExposeHeaderMiddleware(e),
		jwtMiddleware,
	)
	echoApp.HideBanner = true
	return echoApp
}
