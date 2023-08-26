package authenticate

import (
	"bookman/db"
	"crypto/rand"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/sirupsen/logrus"
)

type Auth struct {
	db                    *db.GormDB
	logger                *logrus.Logger
	jwtExpirationDuration time.Duration
	secretKey             []byte
}

func NewAuth(db *db.GormDB, logger *logrus.Logger, jwtExpirationDuration time.Duration) (*Auth, error) {
	secretKey, err := generateRandomKey()
	if err != nil {
		return nil, err
	}

	if db == nil {
		return nil, errors.New("database can not be nil")
	}

	return &Auth{
		db:                    db,
		logger:                logger,
		jwtExpirationDuration: jwtExpirationDuration,
		secretKey:             secretKey,
	}, nil
}

type Credentials struct {
	Username string
	Password string
}

type Token struct {
	TokenString string
}

type claims struct {
	jwt.MapClaims
	Username string `json:"username"`
}

func (a *Auth) Login(cred Credentials) (*Token, error) {

	// Check existence of user
	account, err := a.db.GetUserByUsername(cred.Username)
	if err != nil {
		return nil, err
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(cred.Password))
	if err != nil {
		return nil, errors.New("the password is not correct")
	}

	//	Create JWT token
	expirationTime := time.Now().Add(a.jwtExpirationDuration)
	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		MapClaims: jwt.MapClaims{
			"expired": expirationTime.Unix(),
		},
		Username: cred.Username,
	})

	tokenString, err := tokenJWT.SignedString(a.secretKey)
	if err != nil {
		return nil, err
	}

	return &Token{
		TokenString: tokenString,
	}, nil
}

func (a *Auth) GenerateToken(cred Credentials) (Token, error) {
	return Token{}, nil
}

func (a *Auth) GetAccountByToken(token string) (*string, error) {
	return nil, nil
}

func generateRandomKey() ([]byte, error) {
	jwtKey := make([]byte, 32)
	if _, err := rand.Read(jwtKey); err != nil {
		return nil, err
	}
	return jwtKey, nil
}
