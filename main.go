package main

import (
	"fmt"
	"net/http"

	database "github.com/MJ-NMR/note/database"
	_ "github.com/MJ-NMR/note/modules"
	_ "modernc.org/sqlite"
)


type handler struct {
	db *database.DB
}

func (h *handler) root(w http.ResponseWriter, r *http.Request) {
	notes, err := h.db.GetAllNots()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(401)
		return
	}

	for _, note := range notes {
		fmt.Fprintf(w, "<p>%s</p><p>%s</p><p>%s</p>\n", note.User , note.Content, note.CreatedAt)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	fmt.Println("from the log")
	return next
}


func main() {
	db , err := database.NewDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("connected to the database")

	server := http.NewServeMux()
	handler := &handler{db: db}
	server.HandleFunc("/", handler.root)
	err = http.ListenAndServe(":8000", loggingMiddleware(server))
	if err != nil {
		fmt.Println(err)
		return
	}
}
