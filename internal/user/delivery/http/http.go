package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GianOrtiz/bean/internal/user/usecase"
	"github.com/GianOrtiz/bean/pkg/user"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
)

type UserHTTPHandler struct {
	usecase       user.UseCase
	sessionsStore sessions.Store
}

type registerBody struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func New(usecase user.UseCase, sessionsStore sessions.Store) UserHTTPHandler {
	return UserHTTPHandler{
		usecase:       usecase,
		sessionsStore: sessionsStore,
	}
}

func (h *UserHTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	var body registerBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := validate.Struct(body); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.usecase.Register(body.Email, body.Name, body.Password); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := h.usecase.GetUser(id)
	if err != nil {
		if err == usecase.UserNotFoundErr {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *UserHTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	email, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session, err := h.sessionsStore.Get(r, "uat")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.usecase.Login(email, password, session)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := session.Save(r, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
