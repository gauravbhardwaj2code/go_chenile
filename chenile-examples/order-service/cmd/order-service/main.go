package main

import (
	"log"

	"config"
	"packager"

	"order-service/order/module"
)

func main() {
	cfg, err := config.Load("config/application.yaml")
	if err != nil {
		log.Fatal(err)
	}
	app, err := packager.NewChenileWebApp(module.New())
	if err != nil {
		log.Fatal(err)
	}
	address := ":" + cfg.String("server.port", "8080")
	log.Println("listening on " + address)
	log.Fatal(app.ListenAndServe(address))
}
