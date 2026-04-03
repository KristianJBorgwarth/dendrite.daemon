package repositories

import (
	"database/sql"

	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/store"
)

type UnitOfWork struct {
	dbContext   *sql.DB
	Transaction *sql.Tx
	FileStore   *store.FileStore
}

func NewUnitOfWork() *UnitOfWork {
	db, err := persistence.GetDBContext()
	if err != nil {
		panic("failed to get DB context: " + err.Error())
	}
	return &UnitOfWork{dbContext: db.DB, FileStore: store.NewFileStore()}
}

func (u *UnitOfWork) Begin() (tx *sql.Tx, err error) {
	tx, err = u.dbContext.Begin()
	if err != nil {
		return nil, err
	}
	u.Transaction = tx
	return tx, nil
}

func (u *UnitOfWork) Commit() error {
	if err := u.FileStore.Flush(); err != nil {
		u.Transaction.Rollback()
		u.Transaction = nil
		return err
	}
	if err := u.Transaction.Commit(); err != nil {
		u.FileStore.Rollback()
		u.Transaction = nil
		return err
	}
	u.Transaction = nil
	return nil
}

func (u *UnitOfWork) Rollback() {
	if u.Transaction == nil {
		return
	}
	u.Transaction.Rollback()
	u.FileStore.Rollback()
}
