package main

import (
	"context"
	"log"

	"github.com/beldin0/microservices/pkg/microservice"
)

func main() {
	defer log.Println("Service closed")
	log.Println("Service started")
	s := &service{}
	ms := microservice.New(context.Background())
	err := ms.SetupConnections(
		microservice.MessagebusConnection(nil, &s.messagebus),
	)
	if err != nil {
		log.Fatal(err)
	}
	err = ms.Start(s)
	if err != nil {
		log.Fatal(err)
	}
}
