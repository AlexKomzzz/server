package repository

import (
	chat "github.com/AlexKomzzz/server"
)

// type AuthPostgres struct {
// 	db *sqlx.DB
// }

// func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
// 	return &AuthPostgres{db: db}
// }

// создадим пользователя в БД
// необходимо передать структуру User с зашифрованным паролем
func (r *Repository) CreateUser(user chat.User) (int, error) {

	query := "INSERT INTO (username, email, passwod_hash) VALUES ($1, $2, $3) RETURN id"

	row := r.db.QueryRow(query, user.Username, user.Email, user.Password)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
