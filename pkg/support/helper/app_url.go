package helper

func WithAppUrl(url string) string {

	return ""
	//return fmt.Sprintf(
	//	"%s/%s",
	//	cfg.HTTP.URL,
	//	strings.TrimLeft(url, "/"),
	//)
}

func WithoutAppUrl(url string, noTrailing bool) string {
	return ""
	//cfg, _ := config.NewConfig()
	//
	//str := strings.Replace(
	//	url,
	//	cfg.HTTP.URL,
	//	"",
	//	1,
	//)
	//
	//if noTrailing {
	//	return strings.TrimLeft(str, "/")
	//}
	//
	//return str
}
