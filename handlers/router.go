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
	if !h.authorized(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userId := r.FormValue("id")
	notes, err := h.db.GetAllNots(userId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, note := range notes {
		fmt.Fprintf(w, "<p>%d</p><p>%d</p><p>%s</p><p>%s</p>\n", note.Id, note.User_id, note.Content, note.CreatedAt)
	}
}

func (h *Handler) GetOneNote(w http.ResponseWriter, r *http.Request) {
	if !h.authorized(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userId := r.FormValue("id")
	noteId := r.PathValue("noteId")
	note, err := h.db.GetOneNote(userId, noteId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "<p>%d</p><p>%s</p><p>%s</p>\n", note.User_id, note.Content, note.CreatedAt)
}

func (h *Handler) DeleteOneNote(w http.ResponseWriter, r *http.Request) {
	if !h.authorized(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userId := r.FormValue("userId")
	noteId := r.PathValue("noteId")
	err := h.db.DeleteOneNote(userId, noteId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AddOneNote(w http.ResponseWriter, r *http.Request) {
	if !h.authorized(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	content, userId := r.FormValue("content"), r.FormValue("id")
	if content == "" {
		fmt.Println("no content")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(userId, content)

	err := h.db.AddOneNote(userId, content)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
