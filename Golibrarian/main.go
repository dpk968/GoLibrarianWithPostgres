package main

import (
	"deepak.gupta/GoLibraryAPI/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Create a new router
	router := mux.NewRouter()
	psqlInfo := "host=localhost port=5432 user=postgres " +
        "password=root dbname=library_system sslmode=disable"
	handlers.InitializeDB(psqlInfo)
	defer handlers.CloseDB()
	// Define routes
	router.HandleFunc("/books", handlers.GetAllBooks).Methods("GET")
	router.HandleFunc("/books", handlers.AddBook).Methods("POST")
	router.HandleFunc("/books/{id}", handlers.GetBookByID).Methods("GET")
	router.HandleFunc("/books/{id}", handlers.UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", handlers.DeleteBook).Methods("DELETE")

	// Log requests
	loggedRouter := handlers.LogRequest(router)

	// Start server
	log.Println("Server started on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))
}
