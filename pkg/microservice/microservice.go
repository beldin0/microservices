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

// New returns a new MService with the supplied context
// It also starts up a goroutine with the interrupt listener that cancels the context when an interrupt is received.
func New(ctx context.Context) *MService {
	ctx, cancel := context.WithCancel(ctx)
	go cancelOnInterrupt(cancel)
	m := &MService{ctx: ctx}
	return m
}

// MService collects the functionality to operate a microservice
type MService struct {
	ctx context.Context
}

// SetupConnections creates connections to the supplied Connections concurrently and returns when they have all been created (or all timed out)
func (m *MService) SetupConnections(conns ...*Connection) error {
	if len(conns) == 0 {
		return nil
	}
	return setupConnections(m.ctx, conns)
}

// Start starts the supplied service. It wraps the call to service.Start() with context cancellation and interrupt monitoring
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

func cancelOnInterrupt(cancel context.CancelFunc) {
	defer cancel()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Interrupt signal received")
}
