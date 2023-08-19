package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	client1,_ := NewClient("C1", "E1")
	account1 := NewAccount(client1)
	client2, _ := NewClient("C2", "E2")
	account2 := NewAccount(client2)

	account1.Credit(1000)
	account2.Credit(1000)

	transaction,err := NewTransaction(account1,account2,100)
	assert.Nil(t,err)
	assert.NotNil(t,transaction)
	assert.Equal(t, account2.Balance, 1100.0)
	assert.Equal(t, account1.Balance, 900.0)
}

func TestCreateTransactionInsuficientFunds(t *testing.T) {
	client1,_ := NewClient("C1", "E1")
	account1 := NewAccount(client1)
	client2, _ := NewClient("C2", "E2")
	account2 := NewAccount(client2)

	account1.Credit(1000)
	account2.Credit(1000)

	transaction,err := NewTransaction(account1,account2,2000)
	assert.NotNil(t,err)
	assert.Nil(t,transaction)
	assert.Equal(t, account2.Balance, 1000.0)
	assert.Equal(t, account1.Balance, 1000.0)
	assert.Error(t, err, "insuficient funds")
}