package createaccount

import (
	"testing"

	"github.com/ruhancs/ms-wallet-go/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ClientGatewayMock struct {
	mock.Mock
}

func (m *ClientGatewayMock) Get(id string) (*entity.Client, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Client), args.Error(1)
}

func (m *ClientGatewayMock) Save(client *entity.Client)  error {
	args := m.Called(client)
	return args.Error(0)
}

type AccountGatewayMock struct {
	mock.Mock
}

func (accountGatewayMock *AccountGatewayMock) Save(account *entity.Account) error{
	args := accountGatewayMock.Called(account)//recebe um account
	return args.Error(0)//ou retorna um erro se nao receber account
}

func (accountGatewayMock *AccountGatewayMock) FindById(input string) (*entity.Account, error) {
	args := accountGatewayMock.Called(input)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func (m *AccountGatewayMock) UpdateBalance(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func TestCreateAccountUseCaseExecute(t *testing.T) {
	client, _ := entity.NewClient("C1", "E1")
	clientGatewayMock := &ClientGatewayMock{}
	//quando o metodo get é chamado com o id retorna o client
	clientGatewayMock.On("Get", client.Id).Return(client,nil)

	accountGatewayMock := &AccountGatewayMock{}
	accountGatewayMock.On("Save", mock.Anything).Return(nil)

	usecase := NewCreateAccountUseCase(accountGatewayMock,clientGatewayMock)
	inputDto := CreateAccountInputDto{
		ClientId: client.Id,
	}
	output,err := usecase.Execute(inputDto)

	assert.Nil(t,err)
	assert.NotNil(t,output.Id)
	clientGatewayMock.AssertExpectations(t)
	accountGatewayMock.AssertExpectations(t)
	clientGatewayMock.AssertNumberOfCalls(t, "Get", 1)
	accountGatewayMock.AssertNumberOfCalls(t, "Save", 1)
}