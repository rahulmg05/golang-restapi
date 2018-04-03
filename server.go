package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
)

//Book model
type Book struct {
	ID     string  `json:"id"`
	IBSN   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

//Author struct embedded in Book struct
type Author struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var books []Book

//Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

//Get single book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range books {
		if params["id"] == item.ID {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

//Create new books
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(100000))
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

//Update books
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//Get the body
	var book Book
	json.NewDecoder(r.Body).Decode(&book)
	//Get the url param
	params := mux.Vars(r)
	//update the book
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			book.ID = item.ID
			books = append(books[:index], book)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

//Delete books
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if params["id"] == item.ID {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {
	//Init router
	r := mux.NewRouter()

	books = append(books, Book{ID: "1", IBSN: "7897HHG", Title: "Book-01", Author: &Author{FirstName: "Rahul", LastName: "Marigowda"}})
	books = append(books, Book{ID: "2", IBSN: "78979HG", Title: "Book-02", Author: &Author{FirstName: "Vibha", LastName: "Somashekar"}})
	books = append(books, Book{ID: "3", IBSN: "79979HG", Title: "Book-03", Author: &Author{FirstName: "Random", LastName: "Rahul"}})
	books = append(books, Book{ID: "4", IBSN: "70079HG", Title: "Book-04", Author: &Author{FirstName: "Vibha", LastName: "Random"}})
	books = append(books, Book{ID: "5", IBSN: "71179HG", Title: "Book-05", Author: &Author{FirstName: "HaHa", LastName: "HuHu"}})

	//create route handlers
	r.HandleFunc("/api/v1/books", getBooks).Methods("GET")
	r.HandleFunc("/api/v1/book/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/v1/book", createBook).Methods("POST")
	r.HandleFunc("/api/v1/book/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/v1/book/{id}", deleteBook).Methods("DELETE")

	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r))
}
