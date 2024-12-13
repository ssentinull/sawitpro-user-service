package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/SawitProRecruitment/UserService/model"
	"github.com/golang-jwt/jwt"
)

type AuthInterface interface {
	GenerateJWT(user model.User) (string, error)
}

type Auth struct {
	opt AuthOptions
}

type AuthOptions struct {
	JWTExpiryDuration time.Duration
	PrivateKeyPath    string
}

func InitAuth(opt AuthOptions) AuthInterface {
	auth := Auth{opt: opt}
	return auth
}

func (a Auth) GenerateJWT(user model.User) (string, error) {
	privateKey, err := a.loadPrivateKey()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub": strconv.FormatInt(user.Id, 10),
		"iat": time.Now(),
		"exp": time.Now().Add(a.opt.JWTExpiryDuration),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (a Auth) loadPrivateKey() (*rsa.PrivateKey, error) {
	privateKeyBytes, err := os.ReadFile(a.opt.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	privateKey := (key).(*rsa.PrivateKey)

	return privateKey, nil
}
