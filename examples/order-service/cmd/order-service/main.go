package main

import (
	"log"

	"packager"

	"order-service/order"
)

func main() {
	app, err := packager.NewWebApp(packager.Module{Name: "order", Register: order.Register})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening on :8080")
	log.Fatal(app.ListenAndServe(":8080"))
}
