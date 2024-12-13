package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/SawitProRecruitment/UserService/model"
	"github.com/golang-jwt/jwt"
)

type AuthInterface interface {
	LoadKeys() (*rsa.PrivateKey, *rsa.PublicKey, error)
	GenerateJWT(user model.User) (string, error)
}

type Auth struct {
	opt AuthOptions
}

type AuthOptions struct {
	JWTExpiryDuration time.Duration
	JWTSecretKey      string
}

func InitAuth(opt AuthOptions) AuthInterface {
	auth := Auth{opt: opt}
	return auth
}

func (a Auth) LoadKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	pemString := fmt.Sprintf("-----BEGIN RSA PRIVATE KEY-----\n%s\n-----END RSA PRIVATE KEY-----", a.opt.JWTSecretKey)
	block, _ := pem.Decode([]byte(pemString))
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, err
	}

	privKey := (key).(*rsa.PrivateKey)
	pubKey := &privKey.PublicKey

	return privKey, pubKey, nil
}

func (a Auth) GenerateJWT(user model.User) (string, error) {
	privateKey, _, err := a.LoadKeys()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id":    user.Id,
		"expires_at": time.Now().Add(a.opt.JWTExpiryDuration).UnixMilli(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
