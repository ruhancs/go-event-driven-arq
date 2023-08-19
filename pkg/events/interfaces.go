package events

import (
	"sync"
	"time"
)

type EventInterface interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() interface{}
	SetPayload(payload interface{})
}

//operacoes para executar qnd os eventos sao chamados
//recebe o evento a ser executado
type EventHandlerInterface interface {
	Handle(event EventInterface, waitGroup *sync.WaitGroup)
}

//gerenciador dos eventos registra o evento e dispacha a operacao
type EventDispatcherInterface interface {
	//registra evento no sistema quqndo o evento acontecer executa o handler
	Register(eventName string, handler EventHandlerInterface) error
	//dispara os eventos registrados executa os handlers
	Dispatch(event EventInterface) error
	//remover o evento da fila
	Remove(eventName string, handler EventHandlerInterface) error
	//verificar se existe o evento com o handler
	Has(eventName string, handler EventHandlerInterface) bool
	Clear()
}