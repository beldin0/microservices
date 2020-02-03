package main

import (
	"errors"

	"github.com/beldin0/microservices/pkg/microservice"
)

func (s *service) Connections() []*microservice.Connection {
	if s.connections == nil {
		s.buildConnections()
	}
	return s.connections
}

func (s *service) buildConnections() []*microservice.Connection {
	connections := []*microservice.Connection{}
	connections = append(connections, s.msgbusConn())
	return connections
}

func (s *service) msgbusConn() *microservice.Connection {
	return &microservice.Connection{
		Name: "messagebus",
		Operation: func() error {
			var err error
			s.messagebus, err = func() (interface{}, error) {
				return nil, errors.New("unable to connect")
				// return nil, nil
			}()
			return err
		},
	}
}
