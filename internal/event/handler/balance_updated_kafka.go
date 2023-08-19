package handler

import (
	"fmt"
	"sync"

	"github.com/ruhancs/ms-wallet-go/pkg/events"
	"github.com/ruhancs/ms-wallet-go/pkg/kafka"
)

type UpdateBalanceKafkaHandler struct {
	Kafka *kafka.Producer
}

func NewUpdateBalanceKafkaHandler(kafka *kafka.Producer) *UpdateBalanceKafkaHandler {
	return &UpdateBalanceKafkaHandler{
		Kafka: kafka,
	}
}

func (h *UpdateBalanceKafkaHandler) Handle(message events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	h.Kafka.Publish(message,nil, "balances")//publica no topico balance
	fmt.Println("UpdateBalanceKafkaHandler: ", message.GetPayload())
}