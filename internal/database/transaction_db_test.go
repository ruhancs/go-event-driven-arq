package database

import (
	"database/sql"
	"testing"

	"github.com/ruhancs/ms-wallet-go/internal/entity"
	"github.com/stretchr/testify/suite"
)

type TransactionDBTestSuite struct {
	suite.Suite
	db *sql.DB
	client1 *entity.Client
	client2 *entity.Client
	accountFrom *entity.Account 
	accountTo *entity.Account 
	transactionDB *TransactionDB
}

func (s *TransactionDBTestSuite) SetupSuite() {
	db,err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec("create table clients (id varchar(255), name varchar(255), email varchar(255), created_at date)")
	db.Exec("create table accounts (id varchar(255), client_id varchar(255), balance float, created_at date)")
	db.Exec("create table transactions (id varchar(255), account_id_from varchar(255), account_id_to varchar(255), amount float, created_at date)")
	
	client1,err := entity.NewClient("C1", "E1")
	s.Nil(err)
	s.client1 = client1
	client2,err := entity.NewClient("C2", "E2")
	s.Nil(err)
	s.client2 = client2

	accountFrom := entity.NewAccount(s.client1)
	accountFrom.Balance = 1000.0
	s.accountFrom = accountFrom
	accountTo := entity.NewAccount(s.client2)
	accountTo.Balance = 1000.0
	s.accountTo = accountTo
	
	s.transactionDB = NewTransactionDB(db)
}

func (s *ClientDbTestSuite) TearDownSuiteTransaction() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE clients")
	s.db.Exec("DROP TABLE accounts")
	s.db.Exec("DROP TABLE transactions")
}

func TestAccountDBTestSuiteTransactions(t *testing.T) {
	suite.Run(t, new(TransactionDBTestSuite))
}

func (s *TransactionDBTestSuite) TestCreate() {
	transaction,err := entity.NewTransaction(s.accountFrom,s.accountTo,100)
	s.Nil(err)
	err = s.transactionDB.Create(transaction)
	s.Nil(err)
	
}