package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	Id string
	Name string
	Email string
	Accounts []*Account
	CreatedAt time.Time
	UpdaetdAt time.Time
}

func NewClient(name string, email string) (*Client, error) {
	client :=  &Client{
		Id : uuid.New().String(),
		Name: name,
		Email: email,
		CreatedAt: time.Now(),
		UpdaetdAt: time.Now(),	
	}
	err := client.Validate()
	if err != nil {
		return nil,err
	}
	return client, nil
}

func (c *Client) Validate() error {
	if c.Name == ""{
		return errors.New("name is requerid")
	}
	if c.Email == ""{
		return errors.New("email is requerid")
	}
	return nil
}

func (c *Client) Update(name string, email string) error {
	c.Name = name
	c.Email = email
	c.UpdaetdAt = time.Now()
	err := c.Validate()

	if err != nil {
		return err
	}
	return nil
} 

func (c *Client) AddAcount(account *Account) error {
	if account.Client.Id != c.Id {
		return errors.New("account is already in use with another client")
	}
	c.Accounts = append(c.Accounts, account)
	return nil
}