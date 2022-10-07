package service

import "log"

func (s *Service) Test() {
	err := s.repos.CreateChat(3, 2)
	if err != nil {
		log.Println(err.Error())
	}
}
