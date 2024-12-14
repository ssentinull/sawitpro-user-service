package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/SawitProRecruitment/UserService/model"
	"github.com/golang-jwt/jwt"
)

type AuthInterface interface {
	LoadKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) // REFACTOR: save private & public key in auth struct
	GenerateJWTToken(user model.User) (string, error)
	ValidateJWTToken(tokenStr string) error
	GetUserId(tokenStr string) (int64, error)
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

func (a Auth) GenerateJWTToken(user model.User) (string, error) {
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

func (a Auth) ValidateJWTToken(tokenStr string) error {
	claims, err := a.getClaims(tokenStr)
	if err != nil {
		return err
	}

	expiration := int64(claims["expires_at"].(float64))
	if expiration < time.Now().UnixMilli() {
		return errors.New("token expired")
	}

	return nil
}

func (a Auth) GetUserId(tokenStr string) (int64, error) {
	claims, err := a.getClaims(tokenStr)
	if err != nil {
		return 0, err
	}

	userId := int64(claims["user_id"].(float64))

	return userId, nil
}

func (a Auth) getClaims(tokenStr string) (jwt.MapClaims, error) {
	_, publicKey, err := a.LoadKeys()
	if err != nil {
		return jwt.MapClaims{}, err
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid algorithm")
		}
		return publicKey, nil
	})
	if err != nil {
		return jwt.MapClaims{}, err
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	return claims, nil
}
