package persistence

import (
	"database/sql"
)

type Database struct {
	DB *sql.DB
}

var dbContext *Database


func InitializeDBContext(vaultName string) (error) {
	db, err := InitializeIndex(vaultName)
	if err != nil {
		return err
	}

	dbContext = &Database{DB: db}
	return nil
}

func GetDBContext() (*Database) {
	if dbContext == nil {
		panic("DBContext is not initialized. Call InitializeDbContext first.")
	}
	return dbContext 
}

func CloseDBContext() error {
	if dbContext == nil {
		return nil
	}
	err := dbContext.DB.Close()
	dbContext = nil
	return err
}
