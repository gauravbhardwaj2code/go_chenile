package main

import (
	"log"

	"github.com/ajapro/chenile-go/packager"

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
