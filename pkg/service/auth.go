package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	chat "github.com/AlexKomzzz/server"
	"github.com/dgrijalva/jwt-go"
)

const tokenTTL = 30 * time.Hour

const JWT_SECRET = "rkjk#4#%35FSFJlja#4353KSFjH"
const SOLT = "bt,&#Rkm54FS#$WR2@#nasf!dsfre%"

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(SOLT)))
}

// функция создания Пользователя
// возвращяем id, если пользователь создан
func (service *Service) CreateUser(user chat.User) (int, error) {
	// захешим пароль
	user.Password = generatePasswordHash(user.Password)

	return service.repos.CreateUser(user)
}

func (service *Service) GenerateToken(email, password string) (string, error) {
	// определим id пользователя
	id, err := service.repos.GetUser(email, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{ // генерация токена
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(), // время действия токена
			IssuedAt:  time.Now().Unix(),               //время создания
		},
		id,
	})

	return token.SignedString([]byte(JWT_SECRET))
}

// Парс токена (получаем из токена id)
func (service *Service) ParseToken(accesstoken string) (int, error) {
	token, err := jwt.ParseWithClaims(accesstoken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

// функция получения username по id
func (service *Service) GetUsername(userId int) (string, error) {
	return service.repos.GetUsername(userId)
}
