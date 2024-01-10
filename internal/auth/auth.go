package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	AUTH_TOKEN_NAME   = "uat"
	AUTHORIZATION_KEY = "authorized"
)

type AuthHandler struct {
	sessionsStore sessions.Store
}

func NewHandler(sessionsStore sessions.Store) *AuthHandler {
	return &AuthHandler{
		sessionsStore: sessionsStore,
	}
}

func (h *AuthHandler) Authorize(handler func(w http.ResponseWriter, r *http.Request, session *sessions.Session)) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		session, err := h.sessionsStore.Get(r, AUTH_TOKEN_NAME)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authorized, ok := session.Values["authorized"]
		if !ok || !authorized.(bool) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		handler(w, r, session)
		return
	}
	return http.HandlerFunc(fn)
}
