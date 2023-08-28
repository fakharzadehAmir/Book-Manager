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

type bookRequestResponse struct {
	Name            string       `json:"name"`
	Author          authorInBook `json:"author"`
	Category        string       `json:"category"`
	Volume          uint         `json:"volume"`
	PublishedAt     string       `json:"published_at"`
	Summary         string       `json:"summary"`
	TableOfContents []string     `json:"table_of_contents"`
	Publisher       string       `json:"publisher"`
}

func HandleBooksForPostMethod(w http.ResponseWriter, r *http.Request,
	bm *BookManagerServer, authorizedUser *string) {
	//	Retrieve user from database
	user, err := bm.DB.GetUserByUsername(*authorizedUser)
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

	var br bookRequestResponse
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

func HandleBooksForGetMethod(w http.ResponseWriter, bm *BookManagerServer) {
	//	Get all books of users
	allBooks, err := bm.DB.GetAllBooks()
	if err != nil {
		bm.Logger.WithError(err).Warn("can not retrieve all books")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	// Marshal all books
	var allBooksResponse []bookRequestResponse
	for _, book := range allBooks {

		// get items of table of contents for each book
		contents, err := bm.DB.GetContentsByBookID(book.ID)
		if err != nil {
			bm.Logger.WithError(err).Warn("can not retrieve contents of book ", book.Name)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		//	Get author of each book
		author, err := bm.DB.GetAuthorByID(book.AuthorID)
		if err != nil {
			bm.Logger.WithError(err).Warn("can not retrieve author of book ", book.Name)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		// Add books to list of all books in structure of bookRequestResponse
		allBooksResponse = append(allBooksResponse,
			bookRequestResponse{
				Name: book.Name,
				Author: authorInBook{
					FirstName:   author.FirstName,
					LastName:    author.LastName,
					Nationality: author.Nationality,
					Birthday:    author.Birthday,
				},
				PublishedAt:     book.PublishedAt,
				Publisher:       book.Publisher,
				Summary:         book.Summary,
				Category:        book.Category,
				Volume:          book.Volume,
				TableOfContents: contents,
			})
	}
	response := map[string]interface{}{
		"books": allBooksResponse,
	}

	resBody, _ := json.Marshal(response)
	w.WriteHeader(http.StatusAccepted)
	w.Write(resBody)

}

func (bm *BookManagerServer) HandleBooks(w http.ResponseWriter, r *http.Request) {
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

	// Check Method POST -> add new book, GET -> returns all book
	if r.Method == http.MethodPost {
		HandleBooksForPostMethod(w, r, bm, accountUsername)
	} else if r.Method == http.MethodGet {
		HandleBooksForGetMethod(w, bm)
	}
}
