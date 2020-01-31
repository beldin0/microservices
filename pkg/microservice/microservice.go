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
	Connections() []*Connection
	Start(context.Context) error
	Shutdown(context.Context) error
}

// Run is the method all microservice's main method should call to start the
// service.
func Run(ctx context.Context, service Service) {
	var cancel func()
	ctx, cancel = context.WithCancel(ctx)

	// Set up the interrupt listener
	go cancelOnInterrupt(cancel)

	// Set up the connections that the service requires
	conns := service.Connections()
	err := setupConnections(ctx, conns)
	if err != nil {
		log.Fatal(err)
	}

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
		log.Println(err)
		return
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
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Shutdown took too long, halting.")
			os.Exit(1)
		}
	}
}

func cancelOnInterrupt(cancel context.CancelFunc) {
	defer cancel()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
}
