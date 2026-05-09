package main

import (
	"log"

	mainweb "mainweb-app"
)

func main() {
	app, err := mainweb.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("mainweb app listening on :8080")
	log.Fatal(app.ListenAndServe(":8080"))
}
