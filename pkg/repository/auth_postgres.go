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

	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id"

	row := r.db.QueryRow(query, user.Username, user.Email, user.Password)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) GetUser(email, password string) (int, error) {
	// определим по email и паролю id пользователя
	query := "SELECT id FROM users WHERE email=$1 AND password=$2"
	var id int
	err := r.db.Get(&id, query, email, password)
	if err != nil {
		return 0, err
	}

	return id, nil
}
