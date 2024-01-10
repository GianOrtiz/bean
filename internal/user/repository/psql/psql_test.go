package psql

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/GianOrtiz/bean/pkg/db"
	"github.com/GianOrtiz/bean/pkg/user"
)

const (
	USER_ID    = 1
	USER_EMAIL = "email@email.com"
)

func TestShouldGetUser(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not stablish a new connection with database mock: %v", err)
	}
	defer dbConn.Close()

	db := db.NewSqlDB(dbConn)
	repository := NewPSQLUserRepositoryRepository(db)

	mock.
		ExpectQuery("SELECT id, name, email FROM user").
		WithArgs(USER_ID).
		WillReturnRows(mock.NewRows([]string{"id", "name", "email"}).AddRow(USER_ID, "name", "email"))

	u, err := repository.GetUser(USER_ID)
	if err != nil {
		t.Errorf("expected no error, received: %v", err)
	}

	if u.ID != USER_ID {
		t.Errorf("expected to get user id %d, received %d", u.ID, USER_ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations of executed statements were not met: %v", err)
	}
}

func TestShouldGetPasswordHash(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not stablish a new connection with database mock: %v", err)
	}
	defer dbConn.Close()

	db := db.NewSqlDB(dbConn)
	repository := NewPSQLUserRepositoryRepository(db)

	passwordHash := "hash"
	mock.
		ExpectQuery("SELECT password_hash FROM user").
		WithArgs(USER_ID).
		WillReturnRows(mock.NewRows([]string{"password_hash"}).AddRow(passwordHash))

	retrievedPasswordHash, err := repository.GetPasswordHash(USER_ID)
	if err != nil {
		t.Errorf("expected no error, received: %v", err)
	}

	if retrievedPasswordHash != passwordHash {
		t.Errorf("expected to get passwordHash %q, received %q", passwordHash, retrievedPasswordHash)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations of executed statements were not met: %v", err)
	}
}

func TestShouldCreate(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not stablish a new connection with database mock: %v", err)
	}
	defer dbConn.Close()

	db := db.NewSqlDB(dbConn)
	repository := NewPSQLUserRepositoryRepository(db)

	u := user.User{
		ID:    USER_ID,
		Name:  "name",
		Email: "email",
	}
	passwordHash := "hash"
	mock.
		ExpectPrepare("INSERT INTO user").
		ExpectExec().
		WithArgs(u.Name, u.Email, passwordHash).
		WillReturnResult(sqlmock.NewResult(int64(USER_ID), 1))

	err = repository.Create(&u, passwordHash)
	if err != nil {
		t.Errorf("expected no error, received: %v", err)
	}

	if u.ID != USER_ID {
		t.Errorf("expected to received user id %d, received %d", USER_ID, u.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations of executed statements were not met: %v", err)
	}
}

func TestShouldGetUserByEmail(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not stablish a new connection with database mock: %v", err)
	}
	defer dbConn.Close()

	db := db.NewSqlDB(dbConn)
	repository := NewPSQLUserRepositoryRepository(db)

	mock.
		ExpectQuery("SELECT id, name, email FROM user").
		WithArgs(USER_EMAIL).
		WillReturnRows(mock.NewRows([]string{"id", "name", "email"}).AddRow(USER_ID, "name", USER_EMAIL))

	u, err := repository.GetUserByEmail(USER_EMAIL)
	if err != nil {
		t.Errorf("expected no error, received: %v", err)
	}

	if u.Email != USER_EMAIL {
		t.Errorf("expected to get user email %q, received %q", u.Email, USER_EMAIL)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations of executed statements were not met: %v", err)
	}
}
