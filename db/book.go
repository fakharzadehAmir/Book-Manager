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
	Author          Author `gorm:"foreignKey:AuthorID"`
	CreatedByID     uint
	CreatedBy       User   `gorm:"foreignKey:CreatedByID"`
	Category        string `gorm:"varchar(20)"`
	Volume          uint
	PublishedAt     string
	Summary         string `gorm:"varchar(100)"`
	Publisher       string `gorm:"varchar(20)"`
	TableOfContents []TableOfContent
}

func (gdb *GormDB) CreateNewBook(newBook *Book) error {
	// check duplicate book
	var count int64
	if gdb.db.Model(&Book{}).Where("name = ?", newBook.Name).Count(&count); count > 0 {
		return errors.New("this book is already added")
	}
	return gdb.db.Create(newBook).Error
}
