package gateway

import "github.com/ruhancs/ms-wallet-go/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}