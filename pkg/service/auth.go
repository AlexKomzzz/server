package service

import (
	"crypto/sha1"
	"fmt"
)

const SOLT = "bt,&#Rkm54FS#$WR2@#nasf!dsfre%"

func GeneratePasswordHash(password string) (string, error) {
	hash := sha1.New()
	_, err := hash.Write([]byte(password))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum([]byte(SOLT))), nil
}
