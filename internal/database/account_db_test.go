package database

import (
	"database/sql"
	"testing"

	"github.com/ruhancs/ms-wallet-go/internal/entity"
	"github.com/stretchr/testify/suite"
)

type AccountDBTestSuite struct {
	suite.Suite
	db *sql.DB
	accountDB *AccountDb
	client *entity.Client
}

func (s *AccountDBTestSuite) SetupSuite() {
	db,err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec("create table clients (id varchar(255), name varchar(255), email varchar(255), created_at date)")
	db.Exec("create table accounts (id varchar(255), client_id varchar(255), balance float, created_at date)")
	s.accountDB = NewAccountDB(db)
	s.client,_ = entity.NewClient("C1","E1") 
}

func (s *ClientDbTestSuite) TearDownSuiteAccounts() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE clients")
	s.db.Exec("DROP TABLE accounts")
}

func TestAccountDBTestSuite(t *testing.T) {
	suite.Run(t, new(AccountDBTestSuite))
}

func (s *AccountDBTestSuite) TestSave() {
	account := entity.NewAccount(s.client)
	err := s.accountDB.Save(account)
	s.Nil(err)
}

func (s *AccountDBTestSuite) TestFindById() {
	s.db.Exec(`insert into clients (id,name,email,created_at) values(?,?,?,?)`,
			s.client.Id, s.client.Name, s.client.Email, s.client.CreatedAt)
	account := entity.NewAccount(s.client)
	err := s.accountDB.Save(account)
	s.Nil(err)

	accountFounded,err := s.accountDB.FindById(account.Id)

	s.Nil(err)
	s.NotNil(accountFounded)
	s.Equal(account.Id, accountFounded.Id)
	s.Equal(account.Client.Id, accountFounded.Client.Id)
	s.Equal(account.Client.Name, accountFounded.Client.Name)
	s.Equal(account.Client.Email, accountFounded.Client.Email)
	s.Equal(account.Balance, accountFounded.Balance)
}