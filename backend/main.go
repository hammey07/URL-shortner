package main

import (
	"fmt"
	"net/http"
	"strings"
)

var urlStore = make(map[string]string)
var counter = 100000

func helloHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	originalURL := r.URL.Query().Get("url")
	// originalURL := "facebook.com"
	if originalURL == "" {
		http.Error(w, "Missing ?url=parameter", http.StatusBadRequest)
		return
	}
	// Add protocol if missing
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "https://" + originalURL
	}

	tinyCode := toBase62(uint64(counter))
	counter++

	urlStore[tinyCode] = originalURL

	tinyURL := "https://shrink.fly.dev/" + tinyCode
	fmt.Fprintf(w, "%s", tinyURL)
}

const (
	base         uint64 = 62
	characterSet        = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func toBase62(num uint64) string {
	encoded := ""
	for num > 0 {
		r := num % base
		num /= base
		encoded = string(characterSet[r]) + encoded

	}
	return encoded
}

func main() {
	// helloHandler()
	http.HandleFunc("/shorten", helloHandler)
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
