
package repositories

import (
	"database/sql"
)

type UnitOfWork struct {
	db *sql.DB
	Transaction *sql.Tx
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

func (u *UnitOfWork) Rollback() error {
	if u.Transaction == nil {
		return nil
	}
	return u.Transaction.Rollback()
}

