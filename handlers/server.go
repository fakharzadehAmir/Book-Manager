package handlers

import (
	"bookman/authenticate"
	"bookman/db"

	"github.com/sirupsen/logrus"
)

type BookManagerServer struct {
	DB           *db.GormDB
	Logger       *logrus.Logger
	Authenticate *authenticate.Auth
}
