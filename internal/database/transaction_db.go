package database

import (
	"database/sql"

	"github.com/ruhancs/ms-wallet-go/internal/entity"
)

type TransactionDB struct {
	DB *sql.DB
}

func NewTransactionDB(db *sql.DB) *TransactionDB {
	return &TransactionDB{
		DB: db,
	}
}

func (t *TransactionDB) Create(transaction *entity.Transaction) error {
	stmt,err := t.DB.Prepare(`insert into transactions (id, account_id_from, account_id_to, amount, created_at) values(?,?,?,?,?)`)
	if err != nil {
		return err
	}
	
	defer stmt.Close()
	_,err = stmt.Exec(transaction.Id,transaction.AccountFrom.Id,transaction.AccountTo.Id,transaction.Amount,transaction.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}