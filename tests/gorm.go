package tests

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMockGORM() (testDB *gorm.DB, close func(), mock sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	testDB, err = gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(err)
	}

	return testDB, func() { db.Close() }, mock
}
