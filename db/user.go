package db

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string `gorm:"varchar(25), unique"`
	Firstname   string `gorm:"varchar(25)"`
	Lastname    string `gorm:"varchar(25)"`
	PhoneNumber string `gorm:"varchar(15), unique"`
	Password    string `gorm:"varchar(25)"`
}

func (gdb *GormDB) CreateNewUser(u *User) error {
	// Emcrypting the user password
	if encryptedPW, err := bcrypt.GenerateFromPassword([]byte(u.Password), 4); err != nil {
		return err
	} else {
		u.Password = string(encryptedPW)
	}

	// check duplicate user
	var count int64
	if gdb.db.Model(&User{}).Where("username = ?", u.Username).Count(&count); count > 0 {
		return errors.New("this username is already taken")
	}
	return gdb.db.Create(u).Error

}

func (gdb *GormDB) GetUserByUsername(username string) (*User, error) {
	var user User
	err := gdb.db.Where(&User{Username: username}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
