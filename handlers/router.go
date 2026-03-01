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

func (h *Handler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.db.GetAllNots()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, note := range notes {
		fmt.Fprintf(w, "<p>%d</p><p>%s</p><p>%s</p><p>%s</p>\n", note.Id, note.User, note.Content, note.CreatedAt)
	}
}

func (h *Handler) GetOneNote(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	note, err := h.db.GetOneNote(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "<p>%s</p><p>%s</p><p>%s</p>\n", note.User, note.Content, note.CreatedAt)
}

func (h *Handler) DeleteOneNote(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.db.DeleteOneNote(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AddOneNote(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(401)
		return
	}

	user, exist := r.Form["user"]
	if !exist || user[0] == "" {
		fmt.Println("no name", user)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	content, exist := r.Form["content"]
	if !exist || content[0] == "" {
		fmt.Println("no content")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.db.AddOneNote(user[0], content[0])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
