package main

import (
	"context"

	"github.com/beldin0/microservices/pkg/microservice"
)

var _ microservice.Service = (*service)(nil)

type service struct {
	messagebus context.Context // just using context as an example interface
	running    chan error
}

type bus interface {
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
