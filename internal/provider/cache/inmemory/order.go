package inmemory

import (
	"context"
	"errors"
	"github.com/inspectorvitya/wb-l0/internal/model"
)

func (s *Cache) AddOrder(_ context.Context, order *model.Order) error {
	s.Lock()
	defer s.Unlock()
	_, ok := s.data[order.OrderUID]
	if ok {
		return errors.New("order exist")
	}
	s.data[order.OrderUID] = order
	return nil
}
func (s *Cache) GetByID(_ context.Context, id string) (*model.Order, error) {
	s.RLock()
	defer s.RUnlock()
	val, ok := s.data[id]
	if !ok {
		return nil, model.ErrNotExistOrder
	}
	return val, nil
}
