package server

import (
	"fmt"
	"net/http"

	"github.com/GianOrtiz/bean/internal/auth"
	"github.com/GianOrtiz/bean/internal/config"
	"github.com/GianOrtiz/bean/internal/db"
	userHTTP "github.com/GianOrtiz/bean/internal/user/delivery/http"
	"github.com/GianOrtiz/bean/internal/user/repository/psql"
	userUseCase "github.com/GianOrtiz/bean/internal/user/usecase"
	domainDB "github.com/GianOrtiz/bean/pkg/db"
	"github.com/gorilla/sessions"
)

type Server struct {
	Mux    *http.ServeMux
	Config *config.Config
}

func New() (*Server, error) {
	config, err := config.FromEnv()
	if err != nil {
		return nil, err
	}

	dbConn, err := db.GetDBConnection(*config)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	db := domainDB.NewSqlDB(dbConn)

	sessionStore := sessions.NewFilesystemStore("./temp", []byte{1, 2, 3, 4, 5, 6, 7, 8})

	authHandler := auth.NewHandler(sessionStore)

	userRepository := psql.NewPSQLUserRepositoryRepository(db)
	userUseCase := userUseCase.NewUserUseCase(userRepository)
	userHTTPHandler := userHTTP.New(userUseCase, sessionStore)

	mux.HandleFunc("/login", userHTTPHandler.Login)
	mux.HandleFunc("/register", userHTTPHandler.Register)
	mux.HandleFunc("/users", authHandler.Authorize(userHTTPHandler.GetByID))

	return &Server{
		Mux:    mux,
		Config: config,
	}, nil
}

func (s *Server) Run() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Config.Port), s.Mux)
}
