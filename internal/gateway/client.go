package gateway

import "github.com/ruhancs/ms-wallet-go/internal/entity"

type ClientGateway interface {
	Get(id string) (*entity.Client, error)
	Save(client *entity.Client) error
}