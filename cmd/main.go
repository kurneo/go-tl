package main

import (
	"github.com/kurneo/go-template/internal"
	"github.com/spf13/viper"
	"log"
)

func main() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	viper.SetDefault("APP_HTTP_PORT", "3000")

	app := internal.InitializeApp()
	app.Start()
}
