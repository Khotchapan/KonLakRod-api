package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func ComparePassword(password, hash string) bool {
	log.Println("ComparePassword")
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GeneratePassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
