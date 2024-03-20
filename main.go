package main

import (
	"net/http"

	handler "github.com/SowaOwl/goasync.git/handlers"
)

func main() {
	http.HandleFunc("/", handler.HelloHandle)
	http.ListenAndServe(":8080", nil)
}
