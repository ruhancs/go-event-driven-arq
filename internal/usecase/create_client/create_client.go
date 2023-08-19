package createclient

import (
	"time"

	"github.com/ruhancs/ms-wallet-go/internal/entity"
	"github.com/ruhancs/ms-wallet-go/internal/gateway"
)

type CreateClientInputDto struct {
	Name string
	Email string
}

type CreateClientOutputDto struct {
	Id string
	Name string
	Email string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateClientUseCase struct {
	ClientGateway gateway.ClientGateway 
}

func NewCreateClientUseCase(clientGateway gateway.ClientGateway) *CreateClientUseCase {
	return &CreateClientUseCase{
		ClientGateway: clientGateway,
	}
}

func (u *CreateClientUseCase) Execute(input CreateClientInputDto) (*CreateClientOutputDto, error) {
	client,err := entity.NewClient(input.Name,input.Email)
	if err != nil {
		return nil, err
	}
	
	err = u.ClientGateway.Save(client)
	if err != nil {
		return nil, err
	}

	return &CreateClientOutputDto {
		Id: client.Id,
		Name: client.Name,
		Email: client.Email,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdaetdAt,
	}, nil
}