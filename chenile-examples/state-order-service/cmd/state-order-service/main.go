package main

import (
	"log"

	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/packager"

	"state-order-service/order"
)

func main() {
	app, err := packager.NewWebApp(packager.Module{Name: "state-order", Register: order.Register})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening on :8080")
	log.Fatal(app.ListenAndServe(":8080"))
}
