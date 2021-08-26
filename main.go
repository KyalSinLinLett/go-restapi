package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// book struct
type Book struct {
	ID     string  `json:"id"`
	ISBN   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// author struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// init book var as slice of book struct
var books []Book

// get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// get book by id
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // get params
	// loop thru books and find ID
	for _, book := range books {
		if book.ID == params["id"] {
			json.NewEncoder(w).Encode(book)
			return
		}
	}

	json.NewEncoder(w).Encode(&Book{})
}

// create book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)   // decodes the body and store it in var pointed by &book
	book.ID = strconv.Itoa(rand.Intn(10000000)) // create ID
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

// update book by id
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for idx, book := range books {
		if book.ID == params["id"] {
			books = append(books[:idx], books[idx+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book) // decodes the body and store it in var pointed by &book
			book.ID = params["id"]                    // create ID
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(books)
}

// delete book by id
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for idx, book := range books {
		if book.ID == params["id"] {
			books = append(books[:idx], books[idx+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {
	// init router
	r := mux.NewRouter()

	// mock data - @todo - implement db
	books = append(books, Book{ID: "1", ISBN: "12345A", Title: "Book 1", Author: &Author{Firstname: "John", Lastname: "Smith"}})
	books = append(books, Book{ID: "2", ISBN: "12345B", Title: "Book 2", Author: &Author{Firstname: "Sam", Lastname: "Smith"}})
	books = append(books, Book{ID: "3", ISBN: "12345C", Title: "Book 3", Author: &Author{Firstname: "Larry", Lastname: "Smith"}})
	books = append(books, Book{ID: "4", ISBN: "12345D", Title: "Book 4", Author: &Author{Firstname: "King", Lastname: "Smith"}})
	books = append(books, Book{ID: "5", ISBN: "12345E", Title: "Book 5", Author: &Author{Firstname: "Tom", Lastname: "Smith"}})

	// route handler / endpoints
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	port := ":8000"
	fmt.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(port, r))
}
