package createtransaction

import (
	"context"
	"testing"

	"github.com/ruhancs/ms-wallet-go/internal/entity"
	"github.com/ruhancs/ms-wallet-go/internal/event"
	"github.com/ruhancs/ms-wallet-go/internal/usecase/mocks"
	"github.com/ruhancs/ms-wallet-go/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

type TransactiongatewayMock struct {
	mock.Mock
}

func (gateway *TransactiongatewayMock) Create(transaction *entity.Transaction) error {
	args := gateway.Called(transaction)
	return args.Error(0)
}

func TestCreateTransactionUseCaseExecute(t *testing.T) {
	client1,_ := entity.NewClient("C1","E1")
	accountFrom := entity.NewAccount(client1)
	accountFrom.Credit(1000)
	client2,_ := entity.NewClient("C1","E1")
	accountTo := entity.NewAccount(client2)
	accountTo.Credit(1000)

	mockUnitOfWork := &mocks.UnitOfWorkMock{}
	mockUnitOfWork.On("Do", mock.Anything, mock.Anything).Return(nil)

	

	inputDto := CreateTransactionInputDto{
		AccountIdFrom: accountFrom.Id,
		AccountIdTo: accountTo.Id,
		Amount: 100.0,
	}
	dispatcher := events.NewEventDispatcher()
	eventTransactionCreated := event.NewTransactionCreated()
	eventBalanceUpdated := event.NewBalanceUpdated()
	context := context.Background()

	usecase := NewCreateTransactionUseCase(mockUnitOfWork,dispatcher,eventTransactionCreated, eventBalanceUpdated)
	output,err := usecase.Execute(context,inputDto)

	assert.Nil(t,err)
	assert.NotNil(t,output)
	mockUnitOfWork.AssertExpectations(t)
	mockUnitOfWork.AssertNumberOfCalls(t,"Do",1)
}