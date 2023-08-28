package handlers

import (
	"bookman/db"
	"encoding/json"
	"io"
	"net/http"
)

type authorInBook struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Birthday    string `json:"birthday"`
	Nationality string `json:"nationality"`
}

type newBookRequest struct {
	Name            string       `json:"name"`
	Author          authorInBook `json:"author"`
	Category        string       `json:"category"`
	Volume          uint         `json:"volume"`
	PublishedAt     string       `json:"published_at"`
	Summary         string       `json:"summary"`
	TableOfContents []string     `json:"table_of_contents"`
	Publisher       string       `json:"publisher"`
}

func (bm *BookManagerServer) HandleAddBook(w http.ResponseWriter, r *http.Request) {
	// Check Method
	if r.Method != http.MethodPost {
		bm.Logger.Warn("the add book api is not called by POST method")
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

	// Parse the request body for the new book
	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		bm.Logger.Warn("can not read the body of the request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var br newBookRequest
	err = json.Unmarshal(reqData, &br)
	if err != nil {
		bm.Logger.Warn("can not unmarshal the add book request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Add book with its user which added it
	var contents []db.TableOfContent
	for _, name := range br.TableOfContents {
		contents = append(contents, db.TableOfContent{Item: name})
	}
	err = bm.DB.CreateNewBook(&db.Book{
		Name:        br.Name,
		CreatedBy:   *user,
		Category:    br.Category,
		PublishedAt: br.PublishedAt,
		Publisher:   br.Publisher,
		Summary:     br.Summary,
		Volume:      br.Volume,
		Author: db.Author{
			FirstName:   br.Author.FirstName,
			LastName:    br.Author.LastName,
			Birthday:    br.Author.Birthday,
			Nationality: br.Author.Nationality,
		},
		TableOfContents: contents,
	})
	if err != nil {
		bm.Logger.WithError(err).Warn("can not add new book")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	response := map[string]interface{}{
		"message": "book has been added successfully",
	}

	resBody, _ := json.Marshal(response)
	w.WriteHeader(http.StatusAccepted)
	w.Write(resBody)
}
