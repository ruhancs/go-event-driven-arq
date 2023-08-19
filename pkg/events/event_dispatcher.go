package events

import (
	"errors"
	"fmt"
	"sync"
)

var ErrHandlerAlreadyExist = errors.New("handler already registered")

type EventDispatcher struct {
	//pode ter varios EventHandlerInterface registrados
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func(eventDispatcher *EventDispatcher) Dispatch(event EventInterface) error {
	if handlers, ok := eventDispatcher.handlers[event.GetName()]; ok { //verificar se o evento tem o metodo Handle todos handlers do evento estarao em handlers
		wg := &sync.WaitGroup{}
		for _, handler := range handlers { // fazer loop pelos handlers do evento
			wg.Add(1)
			go handler.Handle(event,wg)// executar o handle em uma thread
		}
		wg.Wait()
	}
	return nil
}

func (eventDispatcher *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	//if retorna se existe o evento com o mesmo nome do eventName
	if _,ok := eventDispatcher.handlers[eventName]; ok {
		//caso exista percorre todos handlers do evento
		for _,h := range eventDispatcher.handlers[eventName] {
			//se um handler ja existir nos eventos retorna o erro
			if h == handler {
				return ErrHandlerAlreadyExist
			}
		}
	}

	//salva o handler no eventDispatcher com o eventName como chave
	eventDispatcher.handlers[eventName] = append(eventDispatcher.handlers[eventName], handler)

	fmt.Println("SALVOU HANDLER")
	fmt.Println(eventDispatcher.handlers[eventName])
	fmt.Println(handler)
	return nil
}

func (eventDispatcher *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	if _, ok := eventDispatcher.handlers[eventName]; ok {
		for i, h := range eventDispatcher.handlers[eventName] {
			if h == handler {
				fmt.Println("ENTROU NO IF PARA REMOVER HANDLER")
				fmt.Println(eventDispatcher.handlers)
				eventDispatcher.handlers[eventName] = append(eventDispatcher.handlers[eventName][:i], eventDispatcher.handlers[eventName][i+1:]...)
				fmt.Println("HANDLER NOVO EVENTO COM HANDLERS")
				fmt.Println(eventDispatcher.handlers[eventName])
				return nil
			}
		}
	}
	return nil
}

func (eventDispatcher *EventDispatcher) Clear() {
	//zera todos eventos
	eventDispatcher.handlers = make(map[string][]EventHandlerInterface)
}

func (eventDispatcher *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	if _,ok := eventDispatcher.handlers[eventName]; ok {
		for _, h := range eventDispatcher.handlers[eventName] {
			if h == handler {
				return true
			}
		}
	}
	return false
}