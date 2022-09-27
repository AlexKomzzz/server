package repository

import (
	"github.com/AlexKomzzz/server/pkg/service"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Id       string `json:"-"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

// создадим пользователя в БД
func (r *AuthPostgres) createUser(user User) (int, error) {

	query := "INSERT INTO (username, email, passwod_hash) VALUES ($1, $2, $3) RETURN id"
	passHash, err := service.GeneratePasswordHash(user.Password)
	if err != nil {
		return 0, err
	}

	row := r.db.QueryRow(query, user.Username, user.Email, passHash)
	var id int
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
