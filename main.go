package main

import (
	"database/server"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/users", server.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/users", server.FetchUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", server.FetchUser).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", server.UpdateUser).Methods(http.MethodPut)

	fmt.Println("server go brr")
	log.Fatal(http.ListenAndServe(":8080", router))

}
