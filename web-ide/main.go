package main

import (
	"log"
	"net/http"
)

// To try out: https://github.com/pressly/chi

func main() {
	// Simple static webserver:
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("./webroot"))))
}
