package main

import (
	"fmt"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"log"
	"net/http"
)

// To try out: https://github.com/pressly/chi

func main() {
	// Simple static webserver:
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.FileServer("/", http.Dir("./webroot"))
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Testing now!")
	})

	r.Get("/test/:id/arrp/:foo", func(w http.ResponseWriter, r *http.Request) {
		// FIXME: Not safe!!
		fmt.Fprintln(w, chi.URLParam(r, "id"))
		fmt.Fprintln(w, chi.URLParam(r, "foo"))
	})
	r.Get("/test/a", func(w http.ResponseWriter, r *http.Request) {
		// FIXME: Not safe!!
		fmt.Fprintln(w, "foo")
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
