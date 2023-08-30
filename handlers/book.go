package handlers

import (
	"bookman/db"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
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

	// Parse the request body for new book
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
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		bm.Logger.Warn("Method is neither POST nor GET")
		return
	}
}

func HandleOneBookForGetMethod(bm *BookManagerServer, w http.ResponseWriter, r *http.Request, bookID uint) {
	book, err := bm.DB.GetABookByID(bookID)
	if err != nil {
		bm.Logger.WithError(err).Warn("can not retrieve book ", bookID)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	// get items of table of contents for each book
	contents, err := bm.DB.GetContentsByBookID(book.ID)
	if err != nil {
		bm.Logger.WithError(err).Warn("can not retrieve contents of book ", book.Name)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	//	Get author of book
	author, err := bm.DB.GetAuthorByID(book.AuthorID)
	if err != nil {
		bm.Logger.WithError(err).Warn("can not retrieve author of book ", book.Name)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	bookResponse := bookRequestResponse{
		Name: book.Name,
		Author: authorInBook{
			FirstName:   author.FirstName,
			LastName:    author.LastName,
			Birthday:    author.Birthday,
			Nationality: author.Nationality,
		},
		Volume:          book.Volume,
		Category:        book.Category,
		Summary:         book.Summary,
		Publisher:       book.Publisher,
		PublishedAt:     book.PublishedAt,
		TableOfContents: contents,
	}

	resBody, err := json.Marshal(bookResponse)
	if err != nil {
		bm.Logger.WithError(err).Warn("can not marshal retrieved book to json", book.Name)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(resBody)
}

func HandleOneBookForDeleteMethod(
	bm *BookManagerServer,
	w http.ResponseWriter,
	r *http.Request,
	loginUsername *string,
	bookID uint) {

	//	Check if login user is the one who created the book with given ID
	usernameBook, err := bm.DB.GetCreatedByUsernameByID(bookID)
	if *usernameBook != *loginUsername {
		bm.Logger.WithError(err).Warn("you didn't add the book with given ID in URL")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	// Check if there is an error for finding the username of given book
	if err != nil {
		bm.Logger.WithError(err).Warn("can not retrieve the book with given ID ")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err = bm.DB.DeleteBookByID(bookID); err != nil {
		bm.Logger.WithError(err).Warn("can not delete the book with given ID ")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	response := map[string]interface{}{
		"message": "book has been deleted successfully",
	}

	resBody, _ := json.Marshal(response)
	w.WriteHeader(http.StatusAccepted)
	w.Write(resBody)
}

func HandleOneBookForPatchMethod(
	bm *BookManagerServer,
	w http.ResponseWriter,
	r *http.Request,
	loginUsername *string,
	bookID uint) {

	//	Check if login user is the one who created the book with given ID
	usernameBook, err := bm.DB.GetCreatedByUsernameByID(bookID)
	if *usernameBook != *loginUsername {
		bm.Logger.WithError(err).Warn("you didn't add the book with given ID in URL")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	// Check if there is an error for finding the username of given book
	if err != nil {
		bm.Logger.WithError(err).Warn("can not retrieve the book with given ID ")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Parse the request body for the book with given ID
	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		bm.Logger.Warn("can not read the body of the request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var br bookRequestResponse
	err = json.Unmarshal(reqData, &br)
	if err != nil {
		bm.Logger.Warn("can not unmarshal the update book request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//update book which login user has added it
	var contents []db.TableOfContent
	for _, name := range br.TableOfContents {
		contents = append(contents, db.TableOfContent{Item: name})
	}
	updatedBook, err := bm.DB.UpdateBookByID(&db.Book{
		Name:        br.Name,
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
	}, bookID)
	if err != nil {
		bm.Logger.WithError(err).Warn("can not add new book")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// get items of table of contents for each book
	updatedContents, err := bm.DB.GetContentsByBookID(updatedBook.ID)
	if err != nil {
		bm.Logger.WithError(err).Warn("can not retrieve contents of book ", updatedBook.Name)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	//	Get author of book
	author, err := bm.DB.GetAuthorByID(updatedBook.AuthorID)
	if err != nil {
		bm.Logger.WithError(err).Warn("can not retrieve author of book ", updatedBook.Name)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	bookResponse := bookRequestResponse{
		Name: updatedBook.Name,
		Author: authorInBook{
			FirstName:   author.FirstName,
			LastName:    author.LastName,
			Birthday:    author.Birthday,
			Nationality: author.Nationality,
		},
		Volume:          updatedBook.Volume,
		Category:        updatedBook.Category,
		Summary:         updatedBook.Summary,
		Publisher:       updatedBook.Publisher,
		PublishedAt:     updatedBook.PublishedAt,
		TableOfContents: updatedContents,
	}

	resBody, err := json.Marshal(bookResponse)
	if err != nil {
		bm.Logger.WithError(err).Warn("can not marshal retrieved book to json", updatedBook.Name)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(resBody)
}

func (bm *BookManagerServer) HandleOneBook(w http.ResponseWriter, r *http.Request) {

	//	Grab Authorization header
	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		bm.Logger.Warn("token empty")
		return
	}

	//	Retrieve the username of user which is login
	loginUsername, err := bm.Authenticate.GetAccountByToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		bm.Logger.WithError(err).Warn("retrieving account: ")
		return
	}

	//	Check value of given id
	pathID := mux.Vars(r)["id"]
	bookID, err := strconv.ParseUint(pathID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bm.Logger.WithError(err).Warn("can not convert id to uint ")
		return
	}

	//	Check Method
	//	DELETE -> delete the given book,
	//	GET -> details of the given book
	//	PUT -> update details of the given book
	if r.Method == http.MethodGet {
		HandleOneBookForGetMethod(bm, w, r, uint(bookID))
	} else if r.Method == http.MethodDelete {
		HandleOneBookForDeleteMethod(bm, w, r, loginUsername, uint(bookID))
	} else if r.Method == http.MethodPatch {
		HandleOneBookForPatchMethod(bm, w, r, loginUsername, uint(bookID))

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		bm.Logger.Warn("Method is not each of PUT, GET or DELETE")
		return
	}

}
