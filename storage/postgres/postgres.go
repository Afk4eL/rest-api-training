package postgres

import (
	"clean-rest-arch/internal/config"
	"clean-rest-arch/internal/models"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Database *gorm.DB
}

func New(config config.Config) (*Database, error) {
	const op = "storage.postgres.New"

	// var dsn string = fmt.Sprintf("host=%s port=%d user=%s "+
	// 	"password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
	// 	config.Host, config.Port, config.User, config.Password, config.DbName)

	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	dbUrl := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//TODO:constants
	if config.Env != "prod" {
		db = db.Debug()
	}

	if err := db.AutoMigrate(&models.UserEntity{}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.AutoMigrate(&models.TaskEntity{}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Database{Database: db}, nil
}

func (db *Database) Stop() error {
	const op = "storage.postgres.Stop"

	storage, err := db.Database.DB()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	storage.Close()

	return nil
}
