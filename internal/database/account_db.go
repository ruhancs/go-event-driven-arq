package database

import (
	"database/sql"

	"github.com/ruhancs/ms-wallet-go/internal/entity"
)

type AccountDb struct {
	DB *sql.DB
}

func NewAccountDB(db *sql.DB) *AccountDb {
	return &AccountDb{
		DB: db,
	}
}

func (a *AccountDb) FindById(id string) (*entity.Account, error) {
	var account entity.Account
	var client entity.Client
	account.Client = &client

	stmt,err := a.DB.Prepare(`select a.id, a.client_id, a.balance, a.created_at,
		c.id, c.name, c.email, c.created_at
		from accounts a inner join clients c on a.client_id=c.id where a.id=?`)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(id)
	err = row.Scan(
		&account.Id,
		&account.Client.Id,
		&account.Balance,
		&account.CreatedAt,
		&client.Id,
		&client.Name,
		&client.Email,
		&client.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &account, nil
} 

func (a *AccountDb) Save(account *entity.Account) error {
	stmt,err := a.DB.Prepare(`insert into accounts (id, client_id, balance, created_at) values (?,?,?,?)`)
	if err != nil {
		return err
	}
	
	defer stmt.Close()
	_,err = stmt.Exec(account.Id,account.Client.Id,account.Balance,account.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (a *AccountDb) UpdateBalance(account *entity.Account) error {
	stmt,err := a.DB.Prepare("update accounts set balance=? where id=?")
	if err != nil {
		return err
	}
	
	defer stmt.Close()
	_, err = stmt.Exec(account.Balance, account.Id)
	if err != nil {
		return err
	}
	return nil
}