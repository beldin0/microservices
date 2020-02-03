package microservice

import "context"

// This file is an example connection that could be required by a microservice

// MessagebusConfig would contain any config requirements for connecting to the messagebus
type MessagebusConfig struct {
}

// MessagebusConnection returns a Connection that can be used to connect to a messagebus service
func MessagebusConnection(config *MessagebusConfig, AtLeastOnce *context.Context) *Connection {
	name := "messagebus"
	return &Connection{
		Name: name,
		Operation: func() error {
			var err error
			*AtLeastOnce, err = func() (context.Context, error) {
				// return nil, errors.New("unable to connect")
				return context.Background(), nil
			}()
			return err
		},
	}
}
