package repository

import "github.com/jmoiron/sqlx"

// type Authorization interface {
// 	createUser()
// }

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}
