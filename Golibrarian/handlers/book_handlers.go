package handlers

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "time"

    "deepak.gupta/GoLibraryAPI/models"
    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

var db *sql.DB

// InitializeDB initializes the PostgreSQL database connection.
func InitializeDB(dataSourceName string) {
    var err error
    db, err = sql.Open("postgres", dataSourceName)
    if err != nil {
        log.Fatal(err)
    }
    // Check if the connection to the database is successful
    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }
}

// LogRequest logs incoming requests.
func LogRequest(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Call the next handler
        next.ServeHTTP(w, r)

        // Log the request
        log.Printf(
            "%s %s %s %s",
            r.RemoteAddr,
            r.Method,
            r.URL.Path,
            time.Since(start),
        )
    })
}

// GetAllBooks returns all books from the database.
func GetAllBooks(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT * FROM books")
    if err != nil {
        http.Error(w, fmt.Sprintf("failed to fetch books: %v", err), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var books []models.Book
    for rows.Next() {
        var book models.Book
        if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year); err != nil {
            http.Error(w, fmt.Sprintf("failed to scan row: %v", err), http.StatusInternalServerError)
            return
        }
        books = append(books, book)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(books)
}

// AddBook adds a new book to the database.
func AddBook(w http.ResponseWriter, r *http.Request) {
    var book models.Book
    if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
        http.Error(w, fmt.Sprintf("failed to decode request body: %v", err), http.StatusBadRequest)
        return
    }

    _, err := db.Exec("INSERT INTO books (title, author, year) VALUES ($1, $2, $3)", book.Title, book.Author, book.Year)
    if err != nil {
        http.Error(w, fmt.Sprintf("failed to insert book: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(book)
}

// GetBookByID retrieves a book by its ID from the database.
func GetBookByID(w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "invalid book ID", http.StatusBadRequest)
        return
    }

    var book models.Book
    err = db.QueryRow("SELECT * FROM books WHERE id = $1", id).Scan(&book.ID, &book.Title, &book.Author, &book.Year)
    if err != nil {
        http.Error(w, fmt.Sprintf("failed to fetch book: %v", err), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(book)
}

// UpdateBook updates a book by its ID in the database.
func UpdateBook(w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "invalid book ID", http.StatusBadRequest)
        return
    }

    var updatedBook models.Book
    if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
        http.Error(w, fmt.Sprintf("failed to decode request body: %v", err), http.StatusBadRequest)
        return
    }

    _, err = db.Exec("UPDATE books SET title = $1, author = $2, year = $3 WHERE id = $4",
        updatedBook.Title, updatedBook.Author, updatedBook.Year, id)
    if err != nil {
        http.Error(w, fmt.Sprintf("failed to update book: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(updatedBook)
}

// DeleteBook deletes a book by its ID from the database.
func DeleteBook(w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "invalid book ID", http.StatusBadRequest)
        return
    }

    _, err = db.Exec("DELETE FROM books WHERE id = $1", id)
    if err != nil {
        http.Error(w, fmt.Sprintf("failed to delete book: %v", err), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// CloseDB closes the database connection.
func CloseDB() {
    db.Close()
}
