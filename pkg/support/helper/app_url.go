package helper

import (
	"fmt"
	"github.com/kurneo/go-template/config"
	"strings"
)

func WithAppUrl(url string) string {
	cfg, _ := config.NewConfig()

	return fmt.Sprintf(
		"%s/%s",
		cfg.HTTP.URL,
		strings.TrimLeft(url, "/"),
	)
}

func WithoutAppUrl(url string, noTrailing bool) string {
	cfg, _ := config.NewConfig()

	str := strings.Replace(
		url,
		cfg.HTTP.URL,
		"",
		1,
	)

	if noTrailing {
		return strings.TrimLeft(str, "/")
	}

	return str
}
