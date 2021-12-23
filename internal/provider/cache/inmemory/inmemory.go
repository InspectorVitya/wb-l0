package inmemory

import (
	"github.com/inspectorvitya/wb-l0/internal/model"
	"sync"
)

type Cache struct {
	sync.RWMutex
	data map[string]*model.Order
}

func New(orders []model.Order) *Cache {
	ordersMap := make(map[string]*model.Order)
	for i := range orders {
		ordersMap[orders[i].OrderUID] = &orders[i]
	}

	return &Cache{
		data: ordersMap,
	}
}
