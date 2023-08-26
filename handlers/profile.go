package handlers

import (
	"encoding/json"
	"net/http"
)

type userInfoResponse struct {
	Username    string `json:"username"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	PhoneNumber string `json:"phone_number"`
}

func (bm *BookManagerServer) HandleProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//	Grab Authorization header
	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		bm.Logger.Warn("token empty")
		return
	}

	//	Retrieve the related account by token
	accountUsername, err := bm.Authenticate.GetAccountByToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		bm.Logger.WithError(err).Warn("retrieving account: ")
		return
	}

	//	Retrieve user from database
	user, err := bm.DB.GetUserByUsername(*accountUsername)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bm.Logger.WithError(err).Warn("retrieving user from db: ")
		return
	}

	//	Create the response body
	res, err := json.Marshal(&userInfoResponse{
		Username:    user.Username,
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		PhoneNumber: user.PhoneNumber,
	})
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
