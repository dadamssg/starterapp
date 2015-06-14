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

type SqlTokenRepository struct {
	db    *sql.DB
	table string
}

func NewPSQLAccessTokenRepository(db *sql.DB) TokenRepository {
	return &SqlTokenRepository{db: db, table: "access_token"}
}

func NewPSQLRefreshTokenRepository(db *sql.DB) TokenRepository {
	return &SqlTokenRepository{db: db, table: "refresh_token"}
}

func (r *SqlTokenRepository) ByToken(token string) (*Token, error) {
	return r.byKey("token", token)
}

func (r *SqlTokenRepository) Add(token *Token) error {
	stmt, err := r.db.Prepare(fmt.Sprintf("INSERT INTO %s(token, user_id, expires_at) VALUES($1,$2,$3)", r.table))
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		token.Token,
		token.UserId,
		token.ExpiresAt)

	return err
}

func (r *SqlTokenRepository) byKey(key string, value interface{}) (*Token, error) {
	token := &Token{}
	query := fmt.Sprintf("SELECT token, user_id, expires_at FROM %s WHERE %s = $1", r.table, key)
	err := r.db.QueryRow(query, value).Scan(
		&token.Token,
		&token.UserId,
		&token.ExpiresAt)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		token.ExpiresAt = token.ExpiresAt.UTC()

		return token, nil
	}
}
