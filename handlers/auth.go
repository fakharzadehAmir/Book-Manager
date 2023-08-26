package handlers

import (
	"bookman/authenticate"
	"bookman/db"
	"encoding/json"
	"io"
	"net/http"
)

type signupRequest struct {
	Username    string `json:"username"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (bm *BookManagerServer) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Check Method
	if r.Method != http.MethodPost {
		bm.Logger.Warn("the signup api is not called by POST method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body for login user
	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		bm.Logger.Warn("can not read the body of the request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var lr loginRequest
	err = json.Unmarshal(reqData, &lr)
	if err != nil {
		bm.Logger.Warn("can not unmarshal the login request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Use authenticate package to validate the credentials
	token, err := bm.Authenticate.Login(authenticate.Credentials{
		Username: lr.Username,
		Password: lr.Password,
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("can not login"))
		return
	}

	response := map[string]interface{}{
		"access_token": token.TokenString,
	}
	resBody, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}

func (bm *BookManagerServer) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	// Check Method
	if r.Method != http.MethodPost {
		bm.Logger.Warn("the signup api is not called by POST method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse the requeste body for the new user
	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		bm.Logger.Warn("can not read the body of the request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var sr signupRequest
	err = json.Unmarshal(reqData, &sr)
	if err != nil {
		bm.Logger.Warn("can not unmarshal the signup request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Add user to the database
	err = bm.DB.CreateNewUser(&db.User{
		Username:    sr.Username,
		Firstname:   sr.Firstname,
		Lastname:    sr.Lastname,
		PhoneNumber: sr.PhoneNumber,
		Password:    sr.Password,
	})
	if err != nil {
		bm.Logger.WithError(err).Warn("can not create new user")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	response := map[string]interface{}{
		"message": "user has been created successfully",
	}

	resBody, _ := json.Marshal(response)
	w.WriteHeader(http.StatusAccepted)
	w.Write(resBody)

}
