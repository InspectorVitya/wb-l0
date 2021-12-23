package streaming

import (
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"strconv"
	"time"
)

func (n *Nats) subscribe(f func(data []byte) error) error {
	ackWait, err := strconv.Atoi(n.cfg.AckWait)
	if err != nil {
		return err
	}

	n.sub, err = n.con.Subscribe(
		n.cfg.Subject,
		func(m *stan.Msg) {
			err := f(m.Data)
			if err != nil {
				log.Printf(err.Error())
			}
			err = m.Ack()
			if err != nil {
				log.Printf(err.Error())
			}
		},
		stan.AckWait(time.Duration(ackWait)*time.Second),
		stan.DurableName(os.Getenv(n.cfg.DurableName)),
		stan.SetManualAckMode())
	if err != nil {
		return err
	}
	return err
}
