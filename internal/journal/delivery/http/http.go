package http

import (
	"encoding/json"
	"net/http"

	journalUseCase "github.com/GianOrtiz/bean/internal/journal/usecase"
	"github.com/GianOrtiz/bean/pkg/journal"
	"github.com/GianOrtiz/bean/pkg/money"
	"github.com/gorilla/sessions"
)

type JournalHandler struct {
	journalUseCase journal.AccountUseCase
}

func NewJournalHandler(journalUseCase journal.AccountUseCase) *JournalHandler {
	return &JournalHandler{
		journalUseCase: journalUseCase,
	}
}

type transactBody struct {
	FromUser     int `json:"from_user"`
	ToUser       int `json:"to_user"`
	ValueAsCents int `json:"value_as_cents"`
}

func (h *JournalHandler) Transact(w http.ResponseWriter, r *http.Request, session *sessions.Session) {
	userID, ok := session.Values["user_id"]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var body transactBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if body.FromUser != userID.(int) && body.ToUser != userID.(int) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := h.journalUseCase.Transact(body.FromUser, body.ToUser, money.FromCents(body.ValueAsCents))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *JournalHandler) FindAccountEntries(w http.ResponseWriter, r *http.Request, session *sessions.Session) {
	userIDStr, ok := session.Values["user_id"]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID := userIDStr.(int)
	account, err := h.journalUseCase.FindUserAccount(userID)
	if err != nil {
		if err == journalUseCase.JournalAccountNotFoundErr {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	entries, err := h.journalUseCase.FindEntries(account.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(entries); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *JournalHandler) FindUserAccount(w http.ResponseWriter, r *http.Request, session *sessions.Session) {
	userIDStr, ok := session.Values["user_id"]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID := userIDStr.(int)
	account, err := h.journalUseCase.FindUserAccount(userID)
	if err != nil {
		if err == journalUseCase.JournalAccountNotFoundErr {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(account); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
