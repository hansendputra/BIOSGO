package main

import (
	"fmt"
	"log"
	"net/http"

	"hansendputra.com/biosgo/config"
	"hansendputra.com/biosgo/routes"
)

func main() {
	db, err := config.ConnectDB()

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := routes.SetupRouter(db)
	appport := "8000"
	fmt.Println("starting web server at http://localhost:" + appport)
	log.Fatal(http.ListenAndServe(":"+appport, router))
}
