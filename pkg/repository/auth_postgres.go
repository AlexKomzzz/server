package repository

import (
	"fmt"

	chat "github.com/AlexKomzzz/server"
)

// созданиепользователя в БД
// необходимо передать структуру User с зашифрованным паролем
func (r *Repository) CreateUser(user chat.User) (int, error) {

	// query := "INSERT INTO users (username, email, password, chats) VALUES ($1, $2, $3, '{0}') RETURNING id"
	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id"

	row := r.db.QueryRow(query, user.Username, user.Email, user.Password)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// определение id пользователя по email и паролю
func (r *Repository) GetUser(email, password string) (int, error) {
	query := "SELECT id FROM users WHERE email=$1 AND password=$2"
	var id int
	err := r.db.Get(&id, query, email, password)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// определение idпользователя по email
func (r *Repository) GetUserByEmail(email string) (int, error) {
	query := "SELECT id FROM users WHERE email=$1"
	var id int
	err := r.db.Get(&id, query, email)
	if err != nil {
		return 0, fmt.Errorf("error: 'select' from GetUserByEmail (repos): нет пользователя с таким email: %v", err)
	}

	return id, nil
}

// определение username пользователя по id
func (r *Repository) GetUsername(userId int) (string, error) {
	query := "SELECT username FROM users WHERE id=$1"
	var username string
	err := r.db.Get(&username, query, userId)
	if err != nil {
		return "", err
	}

	return username, nil
}
