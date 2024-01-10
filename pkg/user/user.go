//go:generate mockgen -source=./user.go -destination=./mock/user.go

package user

import "github.com/gorilla/sessions"

// UserAuthToken is the user authorization token for authentication.
type UserAuthToken string

// User represents an user in the system.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Repository is the data access interface for user data.
type Repository interface {
	// GetUser retrieves a user by its id.
	GetUser(id int) (*User, error)
	// GetUserByEmail retrieves a user by its email.
	GetUserByEmail(email string) (*User, error)
	// Create creates a new user.
	Create(u *User, passwordHash string) error
	// GetPasswordHash retrieves the user password hash by its id.
	GetPasswordHash(id int) (string, error)
}

// UseCase represents the use cases for users.
type UseCase interface {
	// Register registers a new user in the system.
	Register(email, name, password string) error
	// Login logins a user in the system.
	Login(email, password string, session *sessions.Session) error
	// GetUser retrieves the user information.
	GetUser(id int) (*User, error)
}
