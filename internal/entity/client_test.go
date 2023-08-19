package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateClient(t *testing.T) {
	client,err := NewClient("Ju", "email.com")
	if err != nil{
		t.Error("error creating client")
	}

	assert.NotNil(t,client)
	assert.Equal(t,"Ju",client.Name)
	assert.Equal(t,"email.com",client.Email)
}

func TestCreateClientWhenArgsInvalid(t *testing.T) {
	client,err := NewClient("", "")

	assert.NotNil(t,err)
	assert.Nil(t,client)
}

func TestUpdateClient(t *testing.T) {
	client,_ := NewClient("Ju", "email.com")
	err := client.Update("Ma", "newemail.com")
	assert.Nil(t,err)
	assert.Equal(t,client.Name,"Ma")
	assert.Equal(t,client.Email,"newemail.com")

}

func TestUpdateClientWithInvalidArgs(t *testing.T) {
	client,_ := NewClient("Ju", "email.com")
	err := client.Update("", "")
	assert.NotNil(t,err)
}
func TestAddAccountToClient(t *testing.T) {
	client, _ := NewClient("Ju", "email")
	account := NewAccount(client)
	err := client.AddAcount(account)
	assert.Nil(t,err)
	assert.Equal(t, 1,len(client.Accounts))
}