package utils

import "golang.org/x/crypto/bcrypt"

type CryptInterface interface {
	CompareHashAndPassword(hashedPassword, password []byte) error
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
}

type Crypt struct{}

func InitCrypt() CryptInterface {
	c := Crypt{}
	return c
}

func (c Crypt) CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func (c Crypt) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}
