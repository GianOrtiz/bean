package usecase

import (
	"database/sql"
	"testing"

	"github.com/GianOrtiz/bean/pkg/user"
	mock_user "github.com/GianOrtiz/bean/pkg/user/mock"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

const (
	USER_ID       = 1
	USER_EMAIL    = "email@email.com"
	USER_NAME     = "name"
	USER_PASSWORD = "password"
)

func TestShouldGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepository := mock_user.NewMockRepository(ctrl)
	usecase := userUseCase{
		repository: userRepository,
	}

	userRepository.EXPECT().GetUser(USER_ID).Return(&user.User{
		ID:    USER_ID,
		Name:  USER_NAME,
		Email: USER_EMAIL,
	}, nil)

	u, err := usecase.GetUser(USER_ID)
	if err != nil {
		t.Errorf("expected to receive no error on get user, got %v", err)
	}

	if u.ID != USER_ID {
		t.Errorf("expected user with id %d, received id %d", USER_ID, u.ID)
	}
}

func TestShouldGetErrorOnLoginWithWrongCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepository := mock_user.NewMockRepository(ctrl)
	usecase := userUseCase{
		repository: userRepository,
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(USER_PASSWORD), DEFAULT_HASH_COST)

	userRepository.EXPECT().GetUserByEmail(USER_EMAIL).Return(&user.User{
		ID:    USER_ID,
		Name:  USER_NAME,
		Email: USER_EMAIL,
	}, nil)
	userRepository.EXPECT().GetPasswordHash(USER_ID).Return(string(passwordHash), nil)

	err = usecase.Login(USER_EMAIL, "wrong-password", nil)
	if err == nil {
		t.Errorf("expected to receive an error, got %v", err)
	}
}

func TestShouldGetErrorOnLoginWhenUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepository := mock_user.NewMockRepository(ctrl)
	usecase := userUseCase{
		repository: userRepository,
	}

	userRepository.EXPECT().GetUserByEmail(USER_EMAIL).Return(nil, sql.ErrNoRows)

	err := usecase.Login(USER_EMAIL, "wrong-password", nil)
	if err == nil {
		t.Errorf("expected to receive an error, got %v", err)
	}
}

func TestShouldRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepository := mock_user.NewMockRepository(ctrl)
	usecase := userUseCase{
		repository: userRepository,
	}

	userRepository.EXPECT().Create(&user.User{Name: USER_NAME, Email: USER_EMAIL}, gomock.Any()).Return(nil)

	err := usecase.Register(USER_EMAIL, USER_NAME, USER_PASSWORD)
	if err != nil {
		t.Errorf("expected to receive no error, got %v", err)
	}
}
