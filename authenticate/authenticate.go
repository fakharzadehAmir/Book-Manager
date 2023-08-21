package authenticate

import (
	"bookman/db"
	"crypto/rand"
	"errors"
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
	secretKey, err := generateRnadomKey()
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

func (a *Auth) Login(cred Credentials) (Token, error) {
	return Token{}, nil
}

func (a *Auth) GenerateToken(cred Credentials) (Token, error) {
	return Token{}, nil
}

func (a *Auth) GetAccountByToken(token string) (*string, error) {
	return nil, nil
}

func generateRnadomKey() ([]byte, error) {
	jwtKey := make([]byte, 32)
	if _, err := rand.Read(jwtKey); err != nil {
		return nil, err
	}
	return jwtKey, nil
}
