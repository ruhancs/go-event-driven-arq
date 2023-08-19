package createaccount

import (
	"github.com/ruhancs/ms-wallet-go/internal/entity"
	"github.com/ruhancs/ms-wallet-go/internal/gateway"
)


type CreateAccountInputDto struct {
	ClientId string `json:"client_id"`
}

type CreateAccountOutputDto struct {
	Id string
}

type CreateAccountUseCase struct {
	AccountGateway gateway.AccountGateway
	ClientGateway gateway.ClientGateway
}

func NewCreateAccountUseCase( accountGateway gateway.AccountGateway, clientGateway gateway.ClientGateway) *CreateAccountUseCase{
	return &CreateAccountUseCase{
		AccountGateway: accountGateway,
		ClientGateway: clientGateway,
	}
}

func (usecase *CreateAccountUseCase) Execute(input CreateAccountInputDto) (*CreateAccountOutputDto, error) {
	client, err := usecase.ClientGateway.Get(input.ClientId)
	if err != nil {
		return nil, err
	}
	account := entity.NewAccount(client)

	err = usecase.AccountGateway.Save(account)
	if err != nil {
		return nil, err
	}

	return &CreateAccountOutputDto{
		Id: account.Id,
	},nil
}