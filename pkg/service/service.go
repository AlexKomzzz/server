package service

import "github.com/AlexKomzzz/server/pkg/repository"

type Service struct {
	*repository.Repository
}

func NewService(repos *repository.Repository) *Service {
	return &Service{repos}
}
