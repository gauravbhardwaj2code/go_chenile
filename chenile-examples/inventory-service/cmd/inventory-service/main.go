package main

import (
	"log"

	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/config"
	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/packager"

	"inventory-service/inventory/module"
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
