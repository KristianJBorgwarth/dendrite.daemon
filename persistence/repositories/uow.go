package repositories

import (
	"database/sql"

	"github.com/KristianJBorgwarth/dendrite.daemon/persistence/store"
)

type UnitOfWork struct {
	db *sql.DB
	Transaction *sql.Tx
	FileStore *store.FileStore
}

func NewUnitOfWork(db *sql.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (u *UnitOfWork) Begin() (tx *sql.Tx, err error) {
	tx, err = u.db.Begin()
	if err != nil {
		return nil, err
	}
	u.Transaction = tx
	return tx, nil
}

func (u *UnitOfWork) Commit() error {
	if u.Transaction == nil {
		return nil
	}
	return u.Transaction.Commit()
}

func (u *UnitOfWork) Rollback() {
	if u.Transaction == nil {
		return
	}
	u.Transaction.Rollback()
	u.FileStore.Rollback()
}
