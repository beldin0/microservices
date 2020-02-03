package microservice

import "context"

// Service represents the methods to start and stop a microservice
type Service interface {
	Start(context.Context) error
	Shutdown(context.Context) error
}
