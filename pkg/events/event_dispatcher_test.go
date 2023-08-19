package events

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name string
	Payload interface{}
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

func (e *TestEvent) SetPayload(payload interface{}) {
	e.Payload = payload
}

type TestEventHandler struct {
	Id int
}


func (h *TestEventHandler) Handle(event EventInterface, wg *sync.WaitGroup) {

}

type EventDispatcherTestSuite struct {
	suite.Suite
	event1  TestEvent
	event2 TestEvent
	handler1 TestEventHandler
	handler2 TestEventHandler
	handler3 TestEventHandler
	eventDispatcher *EventDispatcher
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.eventDispatcher = NewEventDispatcher()
	suite.handler1 = TestEventHandler{
		Id: 1,
	}
	suite.handler2 = TestEventHandler{
		Id: 2,
	}
	suite.handler3 = TestEventHandler{
		Id: 3,
	}
	suite.event1 = TestEvent{Name: "E1", Payload: "P1"}
	suite.event2 = TestEvent{Name: "E2", Payload: "P2"}

}


func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.Nil(err)
	//verificar a quantidade de eventos com o nome que foi criado acima
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))
	
	err = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler2)
	suite.Nil(err)
	//verificar a quantidade de eventos com o nome que foi criado acima
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	assert.Equal(suite.T(), &suite.handler1, suite.eventDispatcher.handlers[suite.event1.GetName()][0])
	assert.Equal(suite.T(), &suite.handler2, suite.eventDispatcher.handlers[suite.event1.GetName()][1])
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register_WithSameHandler() {
	err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.Nil(err)
	//verificar a quantidade de eventos com o nome que foi criado acima
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))
	
	err = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	
	suite.Equal(ErrHandlerAlreadyExist,err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.Nil(err)
	//verificar a quantidade de eventos com o nome que foi criado acima
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))
	
	err = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler2)
	suite.Nil(err)
	//verificar a quantidade de eventos com o nome que foi criado acima
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))
	
	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.GetName()]))

	suite.eventDispatcher.Clear()
	suite.Equal(0, len(suite.eventDispatcher.handlers))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Has() {
	err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.Nil(err)
	//verificar a quantidade de eventos com o nome que foi criado acima
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))
	
	err = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler2)
	suite.Nil(err)
	//verificar a quantidade de eventos com o nome que foi criado acima
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event1.GetName(), &suite.handler1))
	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event1.GetName(), &suite.handler2))
	assert.False(suite.T(), suite.eventDispatcher.Has(suite.event1.GetName(), &suite.handler3))
}

func (suite *EventDispatcherTestSuite) TestEventDispatch_Remove() {
	err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler1)
	suite.Nil(err)
	//verificar a quantidade de eventos com o nome que foi criado acima
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))
	
	err = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler2)
	suite.Nil(err)
	//verificar a quantidade de eventos com o nome que foi criado acima
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))
	
	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)
	suite.Nil(err)
	//verificar a quantidade de eventos com o nome que foi criado acima
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.GetName()]))

	suite.eventDispatcher.Remove(suite.event1.GetName(), &suite.handler1)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))
	//suite.Equal(suite.T(), &suite.handler2, suite.eventDispatcher.handlers[suite.event1.GetName()][0])//garantir que o handler correto foi removido
	
	suite.eventDispatcher.Remove(suite.event2.GetName(), &suite.handler3)
	suite.Equal(0, len(suite.eventDispatcher.handlers[suite.event2.GetName()]))
}

type MockHandler struct {
	mock.Mock
}

func(m *MockHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
	m.Called(event)
	wg.Done()
}

func (suite *EventDispatcherTestSuite) TestEventDispatch_Dispatch() {
	eventHandler := &MockHandler{}
	eventHandler.On("Handle", &suite.event1).Return()//quando chamar o metodo Handle com o evento nao retorna nada
	suite.eventDispatcher.Register(suite.event1.GetName(), eventHandler)
	suite.eventDispatcher.Dispatch(&suite.event1)//dispara o metodo handle do event1

	eventHandler.AssertExpectations(suite.T())
	eventHandler.AssertNumberOfCalls(suite.T(), "Handle", 1)//verificar que o metodo handle foi chamado
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}