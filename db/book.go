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
