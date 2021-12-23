package application

import (
	"context"
	"encoding/json"
	"github.com/inspectorvitya/wb-l0/internal/model"
	"github.com/inspectorvitya/wb-l0/internal/provider/cache"
	"github.com/inspectorvitya/wb-l0/internal/provider/cache/inmemory"
	"github.com/inspectorvitya/wb-l0/internal/provider/storage"
	"github.com/inspectorvitya/wb-l0/internal/provider/streaming"
)

type App struct {
	db         storage.OrdersStorage
	cache      cache.OrdersCache
	natsStream *streaming.Nats
}

func New(db storage.OrdersStorage, natsStream *streaming.Nats) *App {
	return &App{
		db:         db,
		natsStream: natsStream,
	}
}

func (app *App) Init(ctx context.Context) error {
	orders, err := app.db.GetAll(ctx)
	if err != nil {
		return err
	}
	app.cache = inmemory.New(orders)
	err = app.natsStream.Start(app.newOrderHandler)
	if err != nil {
		return err
	}
	return err
}

func (app *App) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	return app.cache.GetByID(ctx, id)
}

func (app *App) newOrderHandler(data []byte) error {
	newOrder := &model.Order{}
	err := json.Unmarshal(data, &newOrder)
	if err != nil {
		return err
	}
	if newOrder.OrderUID == "" {
		return model.ErrEmptyOrderID
	}
	err = app.db.CreateOrder(context.TODO(), newOrder)
	if err != nil {
		return err
	}
	err = app.cache.AddOrder(context.TODO(), newOrder)
	return err
}
