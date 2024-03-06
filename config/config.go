package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type (
	Config struct {
		App     `yaml:"app"`
		HTTP    `yaml:"http"`
		Log     `yaml:"log"`
		DB      `yaml:"db"`
		JWT     `yaml:"jwt"`
		Storage `yaml:"storage"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
		Env     string `env-required:"true" yaml:"env" env:"APP_ENV"`
		Debug   bool   `env-required:"true" yaml:"debug" env:"APP_DEBUG"`
	}

	JWT struct {
		Secret  string `env-required:"true" yaml:"secret" env:"JWT_SECRET"`
		Timeout string `env-required:"true" yaml:"timeout" env:"JWT_TIMEOUT"`
	}

	HTTP struct {
		Port         string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		ExposeHeader string `env-required:"false" yaml:"expose_header" env:"HTTP_EXPOSE_HEADER"`
		AllowOrigin  string `env-required:"false" yaml:"allow_origin" env:"HTTP_ALLOW_ORIGIN"`
		URL          string `env-required:"true" yaml:"url" env:"HTTP_APP_URL"`
	}

	LogHook    map[string]interface{}
	LogChannel map[string]string
	Log        struct {
		Level    string                `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
		Default  string                `env-required:"true" yaml:"default" env:"LOG_CHANNEL"`
		Channels map[string]LogChannel `env-required:"true" yaml:"channels"`
		Hooks    map[string]LogHook    `env-required:"true" yaml:"hooks"`
	}

	DBConn struct {
		Host          string `env-required:"true" yaml:"host" env:"DB_HOST"`
		Port          string `env-required:"true" yaml:"port" env:"DB_PORT"`
		DatabaseName  string `env-required:"true" yaml:"db_name" env:"DB_NAME"`
		Username      string `env-required:"true" yaml:"db_user" env:"DB_USERNAME"`
		Password      string `env-required:"true" yaml:"db_password" env:"DB_PASSWORD"`
		PoolMax       int    `env-required:"false" yaml:"db_pool_max" env:"DB_POOL_MAX"`
		MaxConnection int    `env-required:"false" yaml:"db_max_connection" env:"DB_MAX_CONNECTION"`
	}
	DB struct {
		Default     string            `env-required:"true" yaml:"default" env:"DB_CONNECTION"`
		Connections map[string]DBConn `env-required:"true" yaml:"connections"`
	}

	DiskCfg map[string]string
	Storage struct {
		Default string             `env-required:"true" yaml:"default" env:"STORAGE_DRIVER"`
		Disks   map[string]DiskCfg `yaml:"disks"`
	}
)

func (h LogHook) Get(key string) interface{} {
	if v, ok := h[key]; ok {
		return v
	}
	return ""
}

func (c LogChannel) Get(key string) string {
	if v, ok := c[key]; ok {
		return v
	}
	return ""
}

func (d DiskCfg) Get(key string) string {
	if v, ok := d[key]; ok {
		return v
	}
	return ""
}

var (
	cfg     *Config
	cfgOnce sync.Once
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	var err error
	cfgOnce.Do(func() {
		config := &Config{}
		err = cleanenv.ReadConfig("config.yml", config)
		if err != nil {
			return
		}

		err = cleanenv.ReadEnv(config)
		if err != nil {
			return
		}

		cfg = config
	})
	return cfg, err
}
