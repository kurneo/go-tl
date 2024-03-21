package helper

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

func WithAppUrl(path string) string {
	return fmt.Sprintf(
		"%s/%s",
		viper.GetString("APP_URL"),
		strings.TrimLeft(path, "/"),
	)
}

func WithoutAppUrl(url string, noTrailing bool) string {
	str := strings.Replace(
		url,
		viper.GetString("APP_URL"),
		"",
		1,
	)

	if noTrailing {
		return strings.TrimLeft(str, "/")
	}

	return str
}

func ResolveOffset(page, limit int) int {
	if page < 0 {
		page = 1
	}
	return (page - 1) * limit
}

func ResolveTotalPages(total int64, limit int) int {

	if total <= 0 {
		return 0
	}

	if limit <= 0 {
		return 0
	}

	totalPages := total / int64(limit)

	if total%int64(limit) > 0 {
		totalPages++
	}

	return int(totalPages)
}
