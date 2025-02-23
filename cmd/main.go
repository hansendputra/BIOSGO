package main

import (
	"fmt"
	"log"
	"net/http"

	"BIOSGO/config"
	"BIOSGO/routes"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := routes.SetupRouter(db)

	fmt.Println("starting web server at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8000", router))
}
