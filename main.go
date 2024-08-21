package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

const (
	dbUser     = "postgres"
	dbPassword = "password"
	dbName     = "url_shortener"
	port       = "8080"
)

var db *sql.DB
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func init() {
	var err error
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 6)
	for i := range result {
		result[i] = charset[rnd.Intn(len(charset))]
	}
	return string(result)
}

func createShortURL(w http.ResponseWriter, r *http.Request) {
	var request map[string]string
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	originalURL, ok := request["url"]
	if !ok {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	shortKey := generateShortKey()
	_, err := db.Exec("INSERT INTO urls (short_key, original_url) VALUES ($1, $2)", shortKey, originalURL)
	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		log.Println("Error inserting into database:", err)
		return
	}

	response := map[string]string{"short_url": fmt.Sprintf("http://localhost:%s/%s", port, shortKey)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[1:]
	var originalURL string
	err := db.QueryRow("SELECT original_url FROM urls WHERE short_key = $1", shortKey).Scan(&originalURL)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/shorten", createShortURL)
	http.HandleFunc("/", redirectToOriginalURL)
	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
