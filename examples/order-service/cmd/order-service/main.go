package main

import (
	"log"

	"github.com/ajapro/chenile-go/packager"

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
