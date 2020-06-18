package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"log"

	// Router add on
	"github.com/gorilla/mux"
	"github.com/chilts/sid"
)

// Books (Model) Struct like class
// string argument it dictating what Field is looking for what is passed to it in this example the body
type Book struct{
	ID string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	Author *Author `json:"author"`
}

type Author struct {
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
}

type Message struct {
	Message string `json:"message"`
}

// Init books var as a  slice Book struct

var books []Book

// Get All Books
func getBooks(w http.ResponseWriter, r *http.Request) {
	// setting a header
	w.Header().Set("Content-Type", "application/json")
	// encoding into JSON (derulo)
	json.NewEncoder(w).Encode(books)
}

// Get one Books
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // will get any params
	// for in loop Essentially
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w,  "Content Not Found", http.StatusNotFound)
}

// Create books
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	// decoding body
	_ = json.NewDecoder(r.Body).Decode(&book)
	for _, item := range books {
		if item.ID == book.ID {
			http.Error(w,  fmt.Sprintf("Item with an ID %s already Exists. If you meant to update one please refer to update book documentation", book.ID), http.StatusBadRequest)
			return
		}
	}
	if book.ID == "" {
		book.ID = sid.IdHex()
	}
	books = append(books, book)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(book)
}

// update Book
func updateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			// slicing books out of slice
			books = append(books[:index], books[index + 1:]...)
			break
		}
	}
	if message.Message != "" {
		createBook(w, r)
	} else {
		http.Error(w,  "Content Not Found", http.StatusNotFound)
	}
}

// delete Book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // will get any params
	for index, item := range books {
		if item.ID == params["id"] {
			// slicing books out of slice
			books = append(books[:index], books[index + 1:]...)
			break
		}
		http.Error(w,  "Content Not Found", http.StatusNotFound)
	}
	w.WriteHeader(http.StatusNoContent)
}

func main()  {
	// Init Router
	r := mux.NewRouter()

	// Mock Data @TODO: implement db
	books = append(books, Book{ID: "1", Isbn:"a123", Title: "My First API", Author: &Author{Firstname:"Kyle", Lastname:"Barnett"}})
	books = append(books, Book{ID: "2", Isbn:"b123", Title: "Creating a RESTful API With Golang", Author: &Author{Firstname:"Elliot", Lastname:"Forbes"}})

	// Route Handlers / Endpoints
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	// {params}
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":3001", r))
}
