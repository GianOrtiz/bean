package psql

import (
	"github.com/GianOrtiz/bean/pkg/db"
	"github.com/GianOrtiz/bean/pkg/user"
)

type psqlUserRepository struct {
	conn db.Queryer
}

func NewPSQLUserRepositoryRepository(db db.DBConn) user.Repository {
	return &psqlUserRepository{conn: db}
}

func (r *psqlUserRepository) GetUser(id int) (*user.User, error) {
	query := `
		SELECT
			id,
			name,
			email
		FROM
			user
		WHERE
			id=$1
	`
	row := r.conn.QueryRow(query, id)
	var u user.User
	err := row.Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *psqlUserRepository) Create(u *user.User, passwordHash string) error {
	query := `
		INSERT INTO
			user(
				name,
				email,
				password_hash
			)
		VALUES(
			$1,
			$2,
			$3
		)
	`
	stmt, err := r.conn.Prepare(query)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(u.Name, u.Email, passwordHash)
	if err != nil {
		return err
	}
	lastInsertedID, _ := res.LastInsertId()
	u.ID = int(lastInsertedID)
	return nil
}

func (r *psqlUserRepository) GetPasswordHash(id int) (string, error) {
	query := `
		SELECT
			password_hash
		FROM
			user
		WHERE
			id=$1
	`
	row := r.conn.QueryRow(query, id)
	var passwordHash string
	if err := row.Scan(&passwordHash); err != nil {
		return "", err
	}
	return passwordHash, nil
}

func (r *psqlUserRepository) GetUserByEmail(email string) (*user.User, error) {
	query := `
		SELECT
			id,
			name,
			email
		FROM
			user
		WHERE
			email=$1
	`
	row := r.conn.QueryRow(query, email)
	var u user.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email); err != nil {
		return nil, err
	}
	return &u, nil
}
