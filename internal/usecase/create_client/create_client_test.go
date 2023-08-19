package createclient

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

func TestCreateClientUseCase_Execute(t *testing.T) {
	clientGatewayMock := &ClientGatewayMock{}
	clientGatewayMock.On("Save", mock.Anything).Return(nil)
	usecase := NewCreateClientUseCase(clientGatewayMock)
	output, err := usecase.Execute(CreateClientInputDto{
		"C1",
		"E1",
	})

	assert.Nil(t,err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.Id)
	assert.Equal(t, "C1", output.Name)
	assert.Equal(t, "E1", output.Email)
	clientGatewayMock.AssertExpectations(t)//garantir que o metodo Save do ClientGateway foi chamado
	clientGatewayMock.AssertNumberOfCalls(t, "Save", 1)
}