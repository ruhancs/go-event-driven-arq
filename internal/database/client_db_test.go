package database

import (
	"database/sql"
	"testing"

	"github.com/ruhancs/ms-wallet-go/internal/entity"
	_"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type ClientDbTestSuite struct {
	suite.Suite
	db *sql.DB
	clientDb *ClientDB
}

func(s *ClientDbTestSuite) SetupSuite() {
	db,err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec("create table clients (id varchar(255), name varchar(255), email varchar(255), created_at date)")
	s.clientDb = NewClientDb(db)
}

func (s *ClientDbTestSuite) TearDownSuiteClients() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE clients")
}

//ativa todos testes atrelados a suite de test
func TestClientDbTestSuite(t *testing.T) {
	suite.Run(t, new(ClientDbTestSuite))
}

func (s *ClientDbTestSuite) TestSaveGet() {
	client,_ := entity.NewClient("C1", "E1")
	s.clientDb.Save(client)

	getClient,err := s.clientDb.Get(client.Id)

	s.Nil(err)
	s.NotNil(getClient)
	s.Equal(client.Id,getClient.Id)
	s.Equal(client.Name,getClient.Name)
	s.Equal(client.Email,getClient.Email)
}