package main

import (
	"bookman/authenticate"
	"bookman/config"
	"bookman/db"
	"bookman/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Read the configuration
	var cfg config.Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err.Error())
	}

	// Create a new instance of logrus
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetReportCaller(true)
	logger.WithField("config", cfg).Infof("Setting up the configuration.")

	// Create a new instance of database
	gormDB, err := db.NewGormDB(cfg)
	if err != nil {
		logger.WithError(err).Fatalln("error in connecting to the database")
	}
	logger.Infoln("connected to the book management database")

	// Create schemas of models
	err = gormDB.CreateSchema()
	if err != nil {
		logger.WithError(err).Fatalln("can not auto migrate")
	}
	logger.Infoln("migrate tables successfully")

	// Create a new instance of authenticate
	auth, err := authenticate.NewAuth(gormDB, logger, 10*time.Minute)
	if err != nil {
		logger.WithError(err).Fatalln("can not create an instance of authenticate")
	}

	bookManagerServer := handlers.BookManagerServer{
		DB:           gormDB,
		Logger:       logger,
		Authenticate: auth,
	}
	router := mux.NewRouter()
	router.HandleFunc("/auth/signup", bookManagerServer.HandleSignUp)
	router.HandleFunc("/auth/login", bookManagerServer.HandleLogin)
	router.HandleFunc("/profile", bookManagerServer.HandleProfile)
	router.HandleFunc("/books", bookManagerServer.HandleBooks)
	router.HandleFunc("/books/{id:[1-9][0-9]*}", bookManagerServer.HandleOneBook)
	http.Handle("/", router)
	logger.WithError(http.ListenAndServe(":8080", nil)).Fatalln("can not run the http server")
}
