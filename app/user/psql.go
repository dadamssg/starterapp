package user

import (
	"database/sql"
	"fmt"
	"strings"
)

type SqlUserRepository struct {
	db *sql.DB
}

func NewPSQLUserRepository(db *sql.DB) UserRepository {
	return &SqlUserRepository{db: db}
}

func (r *SqlUserRepository) ById(id string) (*User, error) {
	return r.byKey("id", id)
}

func (r *SqlUserRepository) ByEmail(email string) (*User, error) {
	return r.byKey("email", canonicalize(email))
}

func (r *SqlUserRepository) ByUsername(username string) (*User, error) {
	return r.byKey("username_canonical", canonicalize(username))
}

func (r *SqlUserRepository) Add(user *User) error {
	stmt, err := r.db.Prepare("INSERT INTO app_user(id,created_at,email,username,username_canonical,password,enabled,confirmation_token) VALUES($1,$2,$3,$4,$5,$6,$7,$8)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		user.Id,
		user.CreatedAt,
		canonicalize(user.Email),
		user.Username,
		canonicalize(user.Username),
		user.Password,
		user.Enabled,
		user.ConfirmationToken)

	return err
}

func (r *SqlUserRepository) byKey(key string, value interface{}) (*User, error) {
	user := &User{}
	query := fmt.Sprintf("SELECT id, created_at, email, username, password, enabled, confirmation_token FROM app_user WHERE %s = $1", key)
	err := r.db.QueryRow(query, value).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Enabled,
		&user.ConfirmationToken)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return user, nil
	}
}

func canonicalize(value string) string {
	return strings.ToLower(value)
}
