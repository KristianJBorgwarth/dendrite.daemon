package persistence

import (
	"database/sql"
)

type DBContext struct {
	DB *sql.DB
}

var dbContext *DBContext


func InitializeDBContext(vaultPath string) (error) {
	db, err := InitializeIndex(vaultPath)
	if err != nil {
		return err
	}

	dbContext = &DBContext{DB: db}
	return nil
}

func GetDBContext() (*DBContext, error) {
	if dbContext == nil {
		panic("DBContext is not initialized. Call InitializeDbContext first.")
	}
	return dbContext, nil
}

func CloseDBContext() error {
	if dbContext == nil {
		return nil
	}
	err := dbContext.DB.Close()
	dbContext = nil
	return err
}
