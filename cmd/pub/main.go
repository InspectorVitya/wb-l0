package main

import (
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/inspectorvitya/wb-l0/internal/model"
	"github.com/nats-io/stan.go"
	uuid "github.com/satori/go.uuid"
	"log"
)

const (
	clusterID = "test-cluster"
	clientID  = "sub-store"
)

func main() {
	nc, err := stan.Connect(clusterID, clientID, stan.NatsURL("nats://127.0.0.1:4223"))
	if err != nil {
		log.Fatalln(err)
	}
	var a model.Order
	_ = faker.SetRandomMapAndSliceSize(3)
	for i := 0; i < 1000; i++ {
		//if i == 100 || i == 600{
		//	_ = nc.Publish("sub", []byte("kak dela?"))
		//	continue
		//}
		err = faker.FakeData(&a)
		if err != nil {
			fmt.Println(err)
		}
		uid := uuid.NewV4()
		a.OrderUID = uid.String()
		bd, _ := json.Marshal(a)
		_ = nc.Publish("sub", bd)
	}

	_ = nc.Close()

}
