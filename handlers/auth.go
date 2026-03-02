package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
)

var sessions = map[string]string{}

func setSession(userId int64) *http.Cookie {
	b := make([]byte, 16)
	rand.Read(b)
	sessionId := hex.EncodeToString(b)
	sessions[sessionId] = fmt.Sprint(userId)

	return &http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		Path:     "/",
		HttpOnly: true,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	name, passowrd := r.FormValue("name"), r.FormValue("password")
	if name == "" || passowrd == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("no name or password")
		return
	}

	userId, err := h.db.AddUser(name, passowrd)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}

	cookie := setSession(userId)
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if h.authorized(r) {
		w.WriteHeader(http.StatusOK)
		return
	}

	name, passowrd := r.FormValue("name"), r.FormValue("password")
	if name == "" || passowrd == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("no name or password")
		return
	}

	userId, err := h.db.GetUser(name, passowrd)
	if err != nil || userId == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("not in the database", err)
		return
	}
	cookie := setSession(userId)
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) authorized(r *http.Request) bool {
	sessionId, err := r.Cookie("session_id")
	if err != nil {
		return false
	}
	_, ok := sessions[sessionId.Value]
	if !ok {
		return false
	}
	fmt.Println("from auth", sessionId)
	err = r.ParseForm()
	if err != nil {
		fmt.Println("FORM err", err)
		return false
	}
	r.Form.Set("id", sessions[sessionId.Value])
	return true
}
