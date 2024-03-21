package main

import (
	"github.com/kurneo/go-template/internal"
	"github.com/spf13/viper"
	"log"
)

func main() {
	viper.SetConfigFile(".env")
	viper.SetDefault("HTTP_PORT", 3000)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := internal.InitializeApp()
	app.Start(viper.GetInt("HTTP_PORT"))
}
