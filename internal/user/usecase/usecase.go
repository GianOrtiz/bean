package usecase

import (
	"github.com/GianOrtiz/bean/internal/auth"
	"github.com/GianOrtiz/bean/pkg/user"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

const DEFAULT_HASH_COST = 8

type userUseCase struct {
	repository user.Repository
}

func NewUserUseCase(repository user.Repository) user.UseCase {
	return &userUseCase{repository: repository}
}

func (r *userUseCase) Register(email, name, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), DEFAULT_HASH_COST)
	if err != nil {
		return err
	}
	u := user.User{
		Name:  name,
		Email: email,
	}
	err = r.repository.Create(&u, string(passwordHash))
	if err != nil {
		return err
	}
	return nil
}

func (uc *userUseCase) Login(email, password string, session *sessions.Session) error {
	u, err := uc.repository.GetUserByEmail(email)
	if err != nil {
		return UserNotFoundErr
	}

	passwordHash, err := uc.repository.GetPasswordHash(u.ID)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return err
	}

	session.Values["user_id"] = u.ID
	session.Values[auth.AUTHORIZATION_KEY] = true
	return nil
}

func (r *userUseCase) GetUser(id int) (*user.User, error) {
	u, err := r.repository.GetUser(id)
	if err != nil {
		return nil, UserNotFoundErr
	}
	return u, nil
}
