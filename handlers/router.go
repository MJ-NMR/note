package handlers

import (
	"fmt"
	"net/http"

	"github.com/MJ-NMR/note/database"
)

type Handler struct {
	db *database.DB
}

func NewHandler(db *database.DB) Handler {
	return Handler{db}
}

func (h *Handler) AllNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.db.GetAllNots()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(401)
		return
	}

	for _, note := range notes {
		fmt.Fprintf(w, "<p>%s</p><p>%s</p><p>%s</p>\n", note.User, note.Content, note.CreatedAt)
	}
}

func (h *Handler) OneNote(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	note, err := h.db.GetOneNote(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(401)
		return
	}

	fmt.Fprintf(w, "<p>%s</p><p>%s</p><p>%s</p>\n", note.User, note.Content, note.CreatedAt)
}
