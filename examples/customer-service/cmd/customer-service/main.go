package main

import (
	"log"

	"packager"

	"customer-service/customer"
)

func main() {
	app, err := packager.NewWebApp(packager.Module{Name: "customer", Register: customer.Register})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening on :8080")
	log.Fatal(app.ListenAndServe(":8080"))
}
