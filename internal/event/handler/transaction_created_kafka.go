package handler

import (
	"fmt"
	"sync"

	"github.com/ruhancs/ms-wallet-go/pkg/events"
	"github.com/ruhancs/ms-wallet-go/pkg/kafka"
)
 

type TransactionCreatedKafkaHandler struct {
	Kafka *kafka.Producer
}

func NewTransactionCreatedKafaka(kafka *kafka.Producer) *TransactionCreatedKafkaHandler {
	return &TransactionCreatedKafkaHandler{
		Kafka: kafka,
	}
}

//publicar a msg no kafka
func(h *TransactionCreatedKafkaHandler) Handle(message events.EventInterface, waitgroup *sync.WaitGroup) {
	defer waitgroup.Done()
	h.Kafka.Publish(message,nil, "transactions")//publica no topico transactions
	fmt.Println("TransactionCreatedKafkaHandler: ", message.GetPayload())
}