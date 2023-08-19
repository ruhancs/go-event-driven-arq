package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	client,_ := NewClient("Ju","email.com")
	account := NewAccount(client)

	assert.NotNil(t,account)
	assert.Equal(t, client.Id, account.Client.Id)
	assert.Equal(t,account.Balance, 0.0)
}

func TestCreateAccountWithNilClient(t *testing.T) {
	account := NewAccount(nil)
	assert.Nil(t,account)
}

func TestCreditAccount(t *testing.T) {
	client,_ := NewClient("Ju","email.com")
	account := NewAccount(client)
	account.Credit(10.5)

	assert.Equal(t,account.Balance,10.5)
}

func TestDebitAccount(t *testing.T) {
	client,_ := NewClient("Ju","email.com")
	account := NewAccount(client)
	account.Credit(100.0)
	account.Debit(50.0)

	assert.Equal(t,account.Balance,50.0)
}