package main

import (
	"context"
	"log"

	"github.com/beldin0/microservices/pkg/microservice"
)

func main() {
	defer log.Println("Service closed")
	log.Println("Service started")
	microservice.Run(context.Background(), &service{})
}

var _ microservice.Service = (*service)(nil)

type service struct {
	connections []*microservice.Connection
	messagebus  interface{}
	running     chan error
}

func (s *service) Start(ctx context.Context) error {
	s.running = make(chan error, 1)
	defer close(s.running)

	// start http / grpc servers here
	// if they are non-blocking, block on context
	// remember to monitor the context for cancellation

	<-ctx.Done()
	return nil
}

func (s *service) Shutdown(ctx context.Context) error {
	return <-s.running
}
