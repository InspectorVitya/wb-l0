package storage

import (
	"context"
	"github.com/inspectorvitya/wb-l0/internal/model"
)

type OrdersStorage interface {
	CreateOrder(ctx context.Context, order *model.Order) error
	GetByID(ctx context.Context, id string) (*model.Order, error)
	GetAll(ctx context.Context) ([]model.Order, error)
}
