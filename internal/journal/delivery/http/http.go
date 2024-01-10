package http

import (
	"encoding/json"
	"net/http"

	journalUseCase "github.com/GianOrtiz/bean/internal/journal/usecase"
	"github.com/GianOrtiz/bean/pkg/journal"
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
	Entries []journal.TransactEntry `json:"entries"`
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

	userIDInEntries := false
	for _, entry := range body.Entries {
		if entry.UserID == userID.(int) {
			userIDInEntries = true
			break
		}
	}

	if !userIDInEntries {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := h.journalUseCase.Transact(body.Entries)
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
