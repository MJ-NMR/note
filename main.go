package main

import (
	"fmt"
	"net/http"

	database "github.com/MJ-NMR/note/database"
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
		fmt.Fprintf(w, "<p>%s</p>\n", note.Content)
	}
}


func main() {
	db , err := database.NewDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	handler := &handler{db: db}
	http.HandleFunc("/", handler.root)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println(err)
		return
	}
}
