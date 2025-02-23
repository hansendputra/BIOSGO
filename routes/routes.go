package routes

import (
	"database/sql"

	"hansendputra.com/biosgo/controllers"
	"hansendputra.com/biosgo/middleware"

	"github.com/gorilla/mux"
)

func SetupRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/login", controllers.Login(db)).Methods("POST")

	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.ValidateToken)

	router.HandleFunc("/users", controllers.GetUsers(db)).Methods("GET")
	router.HandleFunc("/user/{id}", controllers.GetUser(db)).Methods("GET")
	router.HandleFunc("/user", controllers.CreateUser(db)).Methods("POST")
	router.HandleFunc("/user/{id}", controllers.UpdateUser(db)).Methods("PUT")
	router.HandleFunc("/user/{id}", controllers.DeleteUser(db)).Methods("DELETE")
	api.HandleFunc("/peserta", controllers.GetPeserta(db)).Methods("GET")

	return router
}
