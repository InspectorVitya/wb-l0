package cache

import (
	"context"
	"github.com/inspectorvitya/wb-l0/internal/model"
)

type OrdersCache interface {
	AddOrder(ctx context.Context, order *model.Order) error
	GetByID(ctx context.Context, id string) (*model.Order, error)
}
