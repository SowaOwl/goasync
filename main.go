package main

import (
	"net/http"

	handlers "github.com/SowaOwl/goasync.git/http/handlers"
)

func main() {
	http.HandleFunc("/", handlers.WelcomeHandle)
	http.HandleFunc("/api/async", handlers.AsyncHandle)
	http.ListenAndServe(":8080", nil)
}
