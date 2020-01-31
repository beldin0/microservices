package microservice

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
)

func setupConnections(ctx context.Context, conns []*Connection) error {
	var wg sync.WaitGroup
	connectedServices := make(chan struct{}, len(conns))

	for _, conn := range conns {
		conn := conn
		child, _ := context.WithTimeout(ctx, maxConnectionTime)
		wg.Add(1)
		go startConnection(child, &wg, conn, connectedServices)
	}

	wg.Wait()
	if len(connectedServices) != len(conns) {
		return errors.New("unable to connect to all services")
	}
	close(connectedServices)
	return nil
}

func startConnection(ctx context.Context, wg *sync.WaitGroup, conn *Connection, connected chan<- struct{}) {
	defer wg.Done()
	err := errNotConnected
	b := backoff.NewExponentialBackOff()
	for err != nil {
		if ctx.Err() != nil {
			return
		}
		err = backoff.RetryNotify(
			conn.Operation,
			backoff.WithContext(b, ctx),
			notifyWith(conn.Name))
	}
	connected <- struct{}{}
}

func notifyWith(name string) func(error, time.Duration) {
	return func(err error, lapse time.Duration) {
		log.Printf("failed to connect to %s - retrying in %v", name, lapse.Round(time.Second/10).String())
	}
}
