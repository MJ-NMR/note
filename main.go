package main

import (
	"fmt"
	"net/http"

	"github.com/MJ-NMR/note/database"
	handlers "github.com/MJ-NMR/note/handlers"
)

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	db, err := database.NewDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("connected to the database")

	handler := handlers.NewHandler(db)
	router := http.NewServeMux()
	router.HandleFunc("GET /", handler.AllNotes)
	router.HandleFunc("GET /{id}", handler.OneNote)

	fmt.Println("server listening on port 8000")
	err = http.ListenAndServe(":8000", logging(router))
	if err != nil {
		fmt.Println(err)
		return
	}
}
