package databases

import (
	"errors"
	"os"

	"github.com/sigit14ap/simple-dating-app-backend/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DatabaseInitialize() (db *gorm.DB, err error) {
	db, err = connect()
	if err != nil {
		return nil, err
	}

	err = migrate(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getDBConnectionString() string {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	return dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true"
}

func connect() (db *gorm.DB, err error) {
	db, err = gorm.Open(mysql.Open(getDBConnectionString()), &gorm.Config{})
	if err != nil {
		return nil, errors.New("error connecting to database")
	}

	return db, nil
}

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{}, &models.Match{})
	if err != nil {
		return errors.New("error migrating database schema")
	}

	return nil
}
