package db

import (
	"bookman/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Setup Database Connection

type GormDB struct {
	cfg config.Config
	db  gorm.DB
}

func NewGormDB(cfg config.Config) (*GormDB, error) {
	c := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Name,
		cfg.Database.Password)

	// Create a new database connection
	db, err := gorm.Open(postgres.Open(c))
	if err != nil {
		return nil, err
	}
	return &GormDB{
		cfg: cfg,
		db:  *db,
	}, nil
}

func (gdb *GormDB) CreateSchema() error {
	err := gdb.db.AutoMigrate(&User{}, &TableOfContent{}, &Author{}, &Book{})
	if err != nil {
		return err
	}
	return nil
}
