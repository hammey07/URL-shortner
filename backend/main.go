package main

import (
	"fmt"
	"net/http"
	"strconv"
)

var urlStore = make(map[string]string)
var counter = 1000

func helloHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	originalURL := r.URL.Query().Get("url")
	if originalURL == "" {
		http.Error(w, "Missing ?url=parameter", http.StatusBadRequest)
		return
	}

	counter++
	shortCode := strconv.Itoa(counter)

	urlStore[shortCode] = originalURL

	shortURL := "http://localhost:8080/" + shortCode
	fmt.Fprintf(w, "%s", shortURL)
}

func main() {
	http.HandleFunc("/shorten", helloHandler)
	fmt.Print("Server Running on http://localhost:8080")
	http.HandleFunc("/", redirectHandler)
	http.ListenAndServe(":8080", nil)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:]
	originalURL, ok := urlStore[code]
	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
