package application

import (
	"context"
	"encoding/json"
	"github.com/bxcodec/faker/v3"
	"github.com/inspectorvitya/wb-l0/internal/config"
	"github.com/inspectorvitya/wb-l0/internal/model"
	"github.com/inspectorvitya/wb-l0/internal/provider/storage/pgsql"
	"github.com/inspectorvitya/wb-l0/internal/provider/streaming"
	"github.com/nats-io/stan.go"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func startApp() (*App, config.Config, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, cfg, err
	}
	db, err := pgsql.New(cfg.DataBase)
	defer db.Close()
	if err != nil {
		return nil, cfg, err
	}
	ctx := context.Background()
	natsStreaming := streaming.New(cfg.Nats)
	app := New(db, natsStreaming)
	err = app.Init(ctx)
	if err != nil {
		return nil, cfg, err
	}
	return app, cfg, err
}

func TestApp(t *testing.T) {
	var exp *model.Order
	_ = faker.SetRandomMapAndSliceSize(3)
	err := faker.FakeData(&exp)
	require.NoError(t, err)
	cfg, err := config.New()
	require.NoError(t, err)
	db, err := pgsql.New(cfg.DataBase)
	defer db.Close()
	require.NoError(t, err)
	ctx := context.Background()
	natsStreaming := streaming.New(cfg.Nats)
	app := New(db, natsStreaming)
	err = app.Init(ctx)
	require.NoError(t, err)
	t.Run("publish new order and get order", func(t *testing.T) {
		nc, err := stan.Connect(cfg.Nats.ClusterId, cfg.Nats.ClientId+"test", stan.NatsURL(cfg.Nats.URL))
		require.NoError(t, err)
		orderData, err := json.Marshal(exp)
		require.NoError(t, err)
		err = nc.Publish("order", orderData)
		time.Sleep(time.Second)
		require.NoError(t, err)
		actual, err := app.GetOrder(context.Background(), exp.OrderUID)
		require.NoError(t, err)
		require.EqualValues(t, exp, actual)
	})
}
