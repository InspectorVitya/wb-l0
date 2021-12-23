package streaming

import (
	"github.com/inspectorvitya/wb-l0/internal/config"
	"github.com/nats-io/stan.go"

	"log"
)

type Nats struct {
	con stan.Conn
	sub stan.Subscription
	cfg config.Nats
}

func New(cfg config.Nats) *Nats {
	return &Nats{
		cfg: cfg,
	}
}
func (n *Nats) Start(f func(data []byte) error) error {
	conn, err := stan.Connect(n.cfg.ClusterId, n.cfg.ClientId, stan.NatsURL(n.cfg.URL),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Printf("connection lost, reason: %v", reason)
		}))
	if err != nil {
		return err
	}
	n.con = conn
	err = n.subscribe(f)
	if err != nil {
		return err
	}
	return err
}

func (n *Nats) Stop() error {
	err := n.sub.Unsubscribe()
	if err != nil {
		return err
	}
	err = n.con.Close()
	return err
}
