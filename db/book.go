package db

import (
	"errors"
	"gorm.io/gorm"
)

type TableOfContent struct {
	gorm.Model
	BookID uint
	Book   Book `gorm:"foreignKey:BookID"`
	Item   string
}

type Author struct {
	gorm.Model
	FirstName   string `gorm:"varchar(25)"`
	LastName    string `gorm:"varchar(25)"`
	Birthday    string
	Nationality string `gorm:"varchar(25)"`
}

type Book struct {
	gorm.Model
	Name            string `gorm:"varchar(25), unique"`
	AuthorID        uint
	Author          Author `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
	CreatedByID     uint
	CreatedBy       User   `gorm:"foreignKey:CreatedByID"`
	Category        string `gorm:"varchar(20)"`
	Volume          uint
	PublishedAt     string
	Summary         string           `gorm:"varchar(100)"`
	Publisher       string           `gorm:"varchar(20)"`
	TableOfContents []TableOfContent `gorm:"constraint:OnDelete:CASCADE"` // Cascading delete for TableOfContent
}

func (gdb *GormDB) CreateNewBook(newBook *Book) error {
	// check duplicate book
	var count int64
	if gdb.db.Model(&Book{}).Where("name = ?", newBook.Name).Count(&count); count > 0 {
		return errors.New("this book is already added")
	}
	return gdb.db.Create(newBook).Error
}

func (gdb *GormDB) GetCreatedByUsernameByID(bookID uint) (*string, error) {
	book, err := gdb.GetABookByID(bookID)
	if err != nil {
		return nil, err
	}
	username, err := gdb.GetUsernameByID(book.CreatedByID)
	if err != nil {
		return nil, err
	}
	return username, nil
}

func (gdb *GormDB) DeleteBookByID(bookID uint) error {
	// Delete book with given ID if it doesn't exist give its error
	return gdb.db.Delete(&Book{}, bookID).Error
}

func (gdb *GormDB) UpdateBookByID(book *Book, bookID uint) (*Book, error) {
	//	find the book with bookID in postgres database
	var existingBook Book
	err := gdb.db.First(&existingBook, bookID).Error
	if err != nil {
		return nil, err
	}

	//	update given data in our request body
	if book.Name != "" {
		existingBook.Name = book.Name
	}
	if book.Volume != 0 {
		existingBook.Volume = book.Volume
	}
	if book.PublishedAt != "" {
		existingBook.PublishedAt = book.PublishedAt
	}
	if book.Publisher != "" {
		existingBook.Publisher = book.Publisher
	}
	if book.Summary != "" {
		existingBook.Summary = book.Summary
	}
	if book.Category != "" {
		existingBook.Category = book.Category
	}
	if book.TableOfContents != nil {
		err = gdb.db.Where("book_id = ?", bookID).Delete(&TableOfContent{}).Error
		if err != nil {
			return nil, err
		}
		existingBook.TableOfContents = book.TableOfContents
	}
	checkAuthor := Author{
		FirstName:   "",
		LastName:    "",
		Nationality: "",
		Birthday:    "",
	}
	if book.Author != checkAuthor {
		existedAuthor, err := gdb.GetAuthorByID(existingBook.AuthorID)
		if err != nil {
			return nil, err
		}
		// Update the author fields from the request body
		existedAuthor.FirstName = book.Author.FirstName
		existedAuthor.LastName = book.Author.LastName
		existedAuthor.Birthday = book.Author.Birthday
		existedAuthor.Nationality = book.Author.Nationality
		gdb.db.Save(existedAuthor)
	}

	//	save changed data
	err = gdb.db.Save(existingBook).Error
	if err != nil {
		return nil, err
	}

	return &existingBook, nil
}

func (gdb *GormDB) GetAllBooks() ([]Book, error) {
	var allBooks []Book
	err := gdb.db.Find(&allBooks).Error
	if err != nil {
		return nil, err
	}
	return allBooks, nil
}

func (gdb *GormDB) GetABookByID(bookId uint) (*Book, error) {
	var book Book
	err := gdb.db.Where("id = ?", bookId).First(&book).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (gdb *GormDB) GetContentsByBookID(bookID uint) ([]string, error) {
	var contents []string
	err := gdb.db.Model(&TableOfContent{}).Where("book_id = ?", bookID).Pluck("item", &contents).Error
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func (gdb *GormDB) GetAuthorByID(authorID uint) (*Author, error) {
	var author Author
	err := gdb.db.Where("id = ?", authorID).First(&author).Error
	if err != nil {
		return nil, err
	}
	return &author, nil
}
