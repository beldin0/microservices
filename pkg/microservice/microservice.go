package microservice

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func New(ctx context.Context) *MService {
	ctx, cancel := context.WithCancel(ctx)
	go cancelOnInterrupt(cancel)
	m := &MService{ctx: ctx}
	return m
}

type MService struct {
	ctx context.Context
}

func (m *MService) SetupConnections(conns ...*Connection) error {
	if len(conns) == 0 {
		return nil
	}
	return setupConnections(m.ctx, conns)
}

func (m *MService) Start(service Service) error {
	ctx, cancel := context.WithCancel(m.ctx)

	// Start the service
	done := make(chan error, 1)
	go func() {
		done <- service.Start(ctx)
	}()

	// Wait for one of:
	// 1) the call to service.Start to return an error
	// 2) the context to get cancelled
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		log.Println("Cancellation received, shutting down...")
		ctx, cancel = context.WithTimeout(context.Background(), maxShutdownTime)
		go func() {
			err := service.Shutdown(ctx)
			if err != nil {
				log.Println(err)
			}
			cancel()
		}()
		<-ctx.Done()
		switch err := ctx.Err(); err {
		case context.Canceled:
			return nil
		case context.DeadlineExceeded:
			return errors.New("shutdown took too long, halting")
		default:
			return err
		}
	}
}

const maxShutdownTime = 5 * time.Second
const maxConnectionTime = 5 * time.Minute

var errNotConnected = errors.New("service not connected")

// Connection defines the functions that are needed for Backoff to run
// when connecting to an external service
type Connection struct {
	Name      string
	Operation func() error
}

// Service represents the methods to initialise, start and stop a microservice
type Service interface {
	Start(context.Context) error
	Shutdown(context.Context) error
}

func cancelOnInterrupt(cancel context.CancelFunc) {
	defer cancel()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Interrupt signal received")
}
