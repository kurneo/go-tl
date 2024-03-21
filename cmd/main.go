package main

import (
	"fmt"
	"github.com/kurneo/go-template/internal"
	"github.com/spf13/viper"
	"log"
	"reflect"
)

func main() {
	viper.SetConfigFile(".env")
	viper.SetDefault("HTTP_PORT", 3000)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	s := Dispatcher{
		listeners: map[reflect.Type][]Listener{},
	}

	s.Subscribe(func(event Event) {
		fmt.Println("AAA", event)
	}, EventA{})

	s.Subscribe(func(event Event) {
		fmt.Println("BBB", event)
	}, EventB{})

	s.Listen()

	s.Fire(EventA{
		Id: 10,
	})

	s.Fire(EventA{
		Id: 50,
	})

	s.Fire(EventB{
		Name: "Name",
	})

	app := internal.InitializeApp()
	app.Start(viper.GetInt("HTTP_PORT"))
}
