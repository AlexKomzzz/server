package repository

type Authorization interface {
	createUser()
}

type Repository struct {
	Authorization
}

func NewRepository() *Repository {
	return &Repository{}
}
