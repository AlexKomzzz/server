package service

import (
	"crypto/sha1"
	"fmt"

	chat "github.com/AlexKomzzz/server"
)

const SOLT = "bt,&#Rkm54FS#$WR2@#nasf!dsfre%"

func generatePasswordHash(password string) (string, error) {
	hash := sha1.New()
	_, err := hash.Write([]byte(password))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum([]byte(SOLT))), nil
}

// функция создания Пользователя
// возвращяем true, если пользователь создан
func (service *Service) CreateUser(user chat.User) (int, error) {
	var err error
	// захешим пароль
	user.Password, err = generatePasswordHash(user.Password)
	if err != nil {
		return 0, err
	}

	return service.repos.CreateUser(user)
}

func (service *Service) GenerateToken(email, password string) (string, error) {
	return "", nil
}
